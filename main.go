package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	go_fcm_receiver "github.com/morhaviv/go-fcm-receiver"
)

var startTime = time.Now()

type DeviceDetails struct {
	GcmToken      string `json:"gcm_token"`
	FcmToken      string `json:"fcm_token"`
	AndroidId     uint64 `json:"android_id"`
	SecurityToken uint64 `json:"security_token"`
	PrivateKey    string `json:"private_key"`
	AuthSecret    string `json:"auth_secret"`
}

func SaveDeviceDetails(fcmToken, gcmToken string, androidId, securityToken uint64, privateKey, authSecret string) error {
	deviceDetails := DeviceDetails{
		GcmToken:      gcmToken,
		FcmToken:      fcmToken,
		AndroidId:     androidId,
		SecurityToken: securityToken,
		PrivateKey:    privateKey,
		AuthSecret:    authSecret,
	}

	data, err := json.Marshal(deviceDetails)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("device_details.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadDeviceDetails() (*DeviceDetails, error) {
	data, err := ioutil.ReadFile("device_details.json")
	if err != nil {
		return nil, err
	}

	var deviceDetails DeviceDetails
	err = json.Unmarshal(data, &deviceDetails)
	if err != nil {
		return nil, err
	}

	return &deviceDetails, nil
}
func sendWebhookData(webhookUrl string, data []byte) error {
	// Send the data as an HTTP POST request to the webhook URL
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending data to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return fmt.Errorf("webhook request failed with status code: %d", resp.StatusCode)
	}
}

func handleDevice(w http.ResponseWriter, _ *http.Request) {
	deviceDetails, err := LoadDeviceDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deviceDetails)
}

func main() {
	var SenderId int64 = 847572667627
	var webhookUrl string

	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--senderid" && i+1 < len(os.Args) {
			var err error
			SenderId, err = strconv.ParseInt(os.Args[i+1], 10, 64)
			if err != nil {
				fmt.Println("Error parsing SenderId from command line argument, using default value:", err)
			}
		} else if os.Args[i] == "--webhook" && i+1 < len(os.Args) {
			webhookUrl = os.Args[i+1]
		}
	}

	if webhookUrl == "" {
		fmt.Println("Webhook URL is required. Please provide it as a command-line argument: --webhook <url>")
		return
	}
	// start http server
	http.HandleFunc("/device", handleDevice)
	go func() {
		fmt.Println("Starting HTTP server on :8012")
		err := http.ListenAndServe(":8012", nil)
		if err != nil {
			panic(err)
		}
	}()

	// Check if device_details.json file exists
	_, err := os.Stat("device_details.json")
	if os.IsNotExist(err) {
		// Create a new device and save the details
		newDevice := go_fcm_receiver.FCMClient{
			SenderId: SenderId,
			OnDataMessage: func(message []byte) {
				// Ignore the first 10 seconds of data messages
				if time.Since(startTime) < 10*time.Second {
					return
				}
				sendWebhookData(webhookUrl, message)
				fmt.Println("Received a message:", string(message))
			},
		}
		privateKey, authSecret, err := newDevice.CreateNewKeys()
		if err != nil {
			panic(err)
		}
		fcmToken, gcmToken, androidId, securityToken, err := newDevice.Register()
		if err != nil {
			panic(err)
		}
		err = SaveDeviceDetails(fcmToken, gcmToken, androidId, securityToken, privateKey, authSecret)
		if err != nil {
			panic(err)
		}
		fmt.Println("New device details saved to device_details.json")
		fmt.Printf("Listening for messages using new device, senderid: %d\n", SenderId)
		fmt.Println("--TOKEN--: " + fcmToken)
		err = newDevice.StartListening()
		if err != nil {
			panic(err)
		}
	} else {
		// Load the device details from the JSON file
		deviceDetails, err := LoadDeviceDetails()
		if err != nil {
			panic(err)
		}
		// Listen using the loaded device details
		oldDevice := go_fcm_receiver.FCMClient{
			SenderId:      SenderId,
			GcmToken:      deviceDetails.GcmToken,
			FcmToken:      deviceDetails.FcmToken,
			AndroidId:     deviceDetails.AndroidId,
			SecurityToken: deviceDetails.SecurityToken,
			OnDataMessage: func(message []byte) {
				// Ignore the first 10 seconds of data messages
				if time.Since(startTime) < 10*time.Second {
					return
				}
				sendWebhookData(webhookUrl, message)
				fmt.Println("Received a message:", string(message))
			},
		}

		err = oldDevice.LoadKeys(deviceDetails.PrivateKey, deviceDetails.AuthSecret)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Listening for messages using saved device details, senderid: %d\n", SenderId)
		fmt.Println("--TOKEN--: " + deviceDetails.FcmToken)
		err = oldDevice.StartListening()
		if err != nil {
			panic(err)
		}
	}
}

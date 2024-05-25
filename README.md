# FCM Receiver

FCM Receiver is a Go application that can receive FCM (Firebase Cloud Messaging) notifications and forward them as a webhook.

## Features

- Receives FCM notifications from Android devices
- Stores Android device details (FCM token, Android ID, etc.) to a JSON file
- Forwards the FCM notification data to a specified webhook URL
- Provides two HTTP endpoints: `/token` and `/device`
  - `/device` returns the stored device details

## How to Use

### Prerequisites

- Go version 1.16 or newer
- Git

### Building the Application

1. Clone the repository:

   ```
   git clone https://github.com/agusibrahim/fcmreceiver
   ```

2. Navigate to the project directory:

   ```
   cd fcmreceiver
   ```

3. Build the application:

   ```
   make build
   ```

   The command above will generate a binary for the current operating system. If you want to build for a different platform, you can run the following command:

   ```
   GOOS=linux GOARCH=amd64 make build
   ```

   This will generate a binary for the Linux operating system with the amd64 architecture.

### Running the Application

1. Run the application by passing the webhook URL as an argument:

   ```
   ./fcmreceiver --webhook https://example.com/webhook --deviceid 12345678900
   ```

2. The application will start listening for FCM notifications and forward them to the specified webhook URL.

3. You can access the following HTTP endpoints:
   - `/device`: Returns the stored device details

## Contributing

If you find any bugs or have ideas for further development, feel free to open an issue or create a pull request.

## License

This project is licensed under the [MIT License](LICENSE).


The explanation is similar to the previous Indonesian version:

1. The "Features" section describes the main features of the FCM Receiver application.
2. The "How to Use" section provides instructions for building and running the application.
3. The "Contributing" section invites users to report bugs or suggest improvements.
4. The "License" section states the project's license.

Make sure to update the information, links, and other details according to your project.

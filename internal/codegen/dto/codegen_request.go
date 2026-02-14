package dto

// CodeGenRequest holds the parameters needed to generate ESP32 firmware.
type CodeGenRequest struct {
	// IP is the static IP address to assign to the ESP32 device.
	IP string `json:"ip" binding:"required"`

	// HostIP is the backend server IP that the ESP32 will send data to.
	HostIP string `json:"host_ip" binding:"required"`

	// HOSTSSID is the WiFi network name the ESP32 should connect to.
	HOSTSSID string `json:"host_ssid" binding:"required"`

	// HOSTPASS is the WiFi password.
	HOSTPASS string `json:"host_pass" binding:"required"`

	// Port is the backend server port.
	Port int `json:"port"`

	// Protocol is the communication protocol (http/https).
	Protocol string `json:"protocol"`

	// Token is the JWT device authentication token.
	Token string `json:"token" binding:"required"`

	// DeviceName is an optional device identifier.
	DeviceName string `json:"device_name"`

	// BuildTool is the preferred build tool: "platformio" or "arduino-cli".
	// If empty, the system auto-detects the best available tool.
	BuildTool string `json:"build_tool"`

	// Board FQBN for Arduino CLI (e.g., "esp32:esp32:esp32").
	// Only used when build_tool is "arduino-cli".
	BoardFQBN string `json:"board_fqbn"`
}

// CodeGenResponse is the API response after a successful codegen/build.
type CodeGenResponse struct {
	Message     string `json:"message"`
	BuildTool   string `json:"build_tool"`
	BinarySize  int64  `json:"binary_size_bytes,omitempty"`
	BuildID     string `json:"build_id"`
	DownloadURL string `json:"download_url,omitempty"`
}

// UploadRequest holds parameters for OTA firmware upload.
type UploadRequest struct {
	// DeviceIP is the IP address of the target ESP32 device for OTA upload.
	DeviceIP string `json:"device_ip" binding:"required"`

	// BuildID is the ID of a previous build whose firmware should be uploaded.
	// If empty, a fresh build is triggered first.
	BuildID string `json:"build_id"`
}

// UploadResponse is the API response after an OTA upload.
type UploadResponse struct {
	Message  string `json:"message"`
	DeviceIP string `json:"device_ip"`
}

// ToolStatusResponse reports available build tools.
type ToolStatusResponse struct {
	AvailableTools []string `json:"available_tools"`
}

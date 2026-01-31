// management.js - Management pages (devices and users)

const manageDevicesPage = window.location.pathname === "/manage-devices";
if (manageDevicesPage) {
  loadManageDevicesPage();
}

const manageUsersPage = window.location.pathname === "/manage-users";
if (manageUsersPage) {
  loadManageUsersPage();
}

async function loadManageDevicesPage() {
  const devicesTableBody = document.getElementById("devicesTableBody");
  const addDeviceBtn = document.getElementById("addDeviceBtn");

  if (!devicesTableBody) return;

  try {
    const res = await fetch("/api/devices");
    const data = await res.json();

    if (!res.ok) {
      devicesTableBody.innerHTML = "<tr><td colspan=\"6\">Failed to load devices</td></tr>";
      return;
    }

    renderDevicesTable(data.devices || []);
  } catch (error) {
    devicesTableBody.innerHTML = "<tr><td colspan=\"6\">Failed to load devices</td></tr>";
  }

  // Add device button
  if (addDeviceBtn) {
    addDeviceBtn.addEventListener("click", () => {
      showDeviceModal();
    });
  }

  // Device form
  const deviceForm = document.getElementById("deviceForm");
  if (deviceForm) {
    deviceForm.addEventListener("submit", handleDeviceSubmit);
  }

  // Modal controls
  const closeModal = document.getElementById("closeModal");
  if (closeModal) {
    closeModal.addEventListener("click", () => {
      document.getElementById("deviceModal").style.display = "none";
    });
  }

  // Load device types
  loadDeviceTypes();
  loadVersions();
}

function renderDevicesTable(devices) {
  const devicesTableBody = document.getElementById("devicesTableBody");

  if (!devices.length) {
    devicesTableBody.innerHTML = "<tr><td colspan=\"6\">No devices found</td></tr>";
    return;
  }

  devicesTableBody.innerHTML = devices.map(device => `
    <tr>
      <td>${device.id}</td>
      <td>${device.name}</td>
      <td>${device.type}</td>
      <td>${device.ip_address || 'N/A'}</td>
      <td>${device.mac_address || 'N/A'}</td>
      <td>
        <button class="button small" onclick="editDevice(${device.id})">Edit</button>
        <button class="button small danger" onclick="deleteDevice(${device.id})">Delete</button>
      </td>
    </tr>
  `).join('');
}

async function loadDeviceTypes() {
  try {
    const res = await fetch("/api/device-types");
    const data = await res.json();

    if (res.ok) {
      let deviceTypes = data.device_types || [];
      // Fallback if no device types from API
      if (deviceTypes.length === 0) {
        deviceTypes = [
          {id: 1, name: "volt-current-meter"},
          {id: 2, name: "smart-switch"},
          {id: 3, name: "sensor-node"},
          {id: 4, name: "temperature-sensor"},
          {id: 5, name: "humidity-sensor"},
          {id: 6, name: "motion-detector"},
          {id: 7, name: "relay-module"},
          {id: 8, name: "power-monitor"},
          {id: 9, name: "energy-meter"}
        ];
      }
      const deviceTypeSelect = document.getElementById("deviceType");
      if (deviceTypeSelect) {
        deviceTypeSelect.innerHTML = '<option value="">Select device type</option>' +
          deviceTypes.map(type => `<option value="${type.id}">${type.name}</option>`).join('');
      }
    }
  } catch (error) {
    console.error("Failed to load device types:", error);
    // Fallback on error
    const deviceTypeSelect = document.getElementById("deviceType");
    if (deviceTypeSelect) {
      deviceTypeSelect.innerHTML = '<option value="">Select device type</option>' +
        '<option value="1">volt-current-meter</option>' +
        '<option value="2">smart-switch</option>' +
        '<option value="3">sensor-node</option>' +
        '<option value="4">temperature-sensor</option>' +
        '<option value="5">humidity-sensor</option>' +
        '<option value="6">motion-detector</option>' +
        '<option value="7">relay-module</option>' +
        '<option value="8">power-monitor</option>' +
        '<option value="9">energy-meter</option>';
    }
  }
}

async function loadVersions() {
  try {
    const res = await fetch("/api/versions");
    const data = await res.json();

    if (res.ok) {
      const versions = data.versions || [];
      const deviceFirmwareSelect = document.getElementById("deviceFirmware");
      if (deviceFirmwareSelect) {
        deviceFirmwareSelect.innerHTML = '<option value="">Select firmware version</option>' +
          versions.map(version => `<option value="${version.ID}">${version.Version}</option>`).join('');
      }
    }
  } catch (error) {
    console.error("Failed to load versions:", error);
  }
}

function showDeviceModal(device = null) {
  const modal = document.getElementById("deviceModal");
  const form = document.getElementById("deviceForm");
  const modalTitle = document.getElementById("modalTitle");

  if (device) {
    modalTitle.textContent = "Edit Device";
    document.getElementById("deviceId").value = device.id;
    document.getElementById("deviceName").value = device.name;
    document.getElementById("deviceType").value = device.type;
    document.getElementById("deviceIP").value = device.ip_address || '';
    document.getElementById("deviceMAC").value = device.mac_address || '';
    document.getElementById("deviceFirmware").value = device.version_id || '';
    document.getElementById("deviceAddress").value = device.address || '';
    document.getElementById("deviceCity").value = device.city || '';
  } else {
    modalTitle.textContent = "Add Device";
    form.reset();
    document.getElementById("deviceId").value = '';
  }

  modal.style.display = "block";
}

async function editDevice(id) {
  try {
    const res = await fetch(`/api/devices/${id}`);
    const device = await res.json();

    if (res.ok) {
      showDeviceModal(device);
    }
  } catch (error) {
    alert("Failed to load device details");
  }
}

async function deleteDevice(id) {
  if (!confirm("Are you sure you want to delete this device?")) return;

  try {
    const res = await fetch(`/api/devices/${id}`, {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${getToken()}`
      }
    });

    if (res.ok) {
      loadManageDevicesPage(); // Reload the page
    } else {
      alert("Failed to delete device");
    }
  } catch (error) {
    alert("Failed to delete device");
  }
}

async function handleDeviceSubmit(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const deviceData = {
    name: formData.get("name"),
    type: parseInt(formData.get("type")) || 0,
    ip_address: formData.get("ip_address"),
    mac_address: formData.get("mac_address"),
    firmware_version_id: parseInt(formData.get("firmware_version_id")) || 0,
    address: formData.get("address"),
    city: formData.get("city")
  };

  const deviceId = formData.get("deviceId");
  const isEdit = deviceId && deviceId !== '';

  try {
    const res = await fetch(isEdit ? `/api/devices/${deviceId}` : "/api/devices", {
      method: isEdit ? "PUT" : "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getToken()}`
      },
      body: JSON.stringify(deviceData)
    });

    if (res.ok) {
      document.getElementById("deviceModal").style.display = "none";
      loadManageDevicesPage(); // Reload the page
    } else {
      const error = await res.json();
      alert(error.error || "Failed to save device");
    }
  } catch (error) {
    alert("Failed to save device");
  }
}

async function loadManageUsersPage() {
  const usersTableBody = document.getElementById("usersTableBody");
  const addUserBtn = document.getElementById("addUserBtn");

  if (!usersTableBody) return;

  try {
    const res = await fetch("/api/users", {
      headers: {
        "Authorization": `Bearer ${getToken()}`
      }
    });
    const data = await res.json();

    if (!res.ok) {
      usersTableBody.innerHTML = "<tr><td colspan=\"4\">Failed to load users</td></tr>";
      return;
    }

    renderUsersTable(data || []);
  } catch (error) {
    usersTableBody.innerHTML = "<tr><td colspan=\"4\">Failed to load users</td></tr>";
  }

  // Add user button
  if (addUserBtn) {
    addUserBtn.addEventListener("click", () => {
      showUserModal();
    });
  }

  // User form
  const userForm = document.getElementById("userForm");
  if (userForm) {
    userForm.addEventListener("submit", handleUserSubmit);
  }

  // Modal controls
  const closeUserModal = document.getElementById("closeUserModal");
  if (closeUserModal) {
    closeUserModal.addEventListener("click", () => {
      document.getElementById("userModal").style.display = "none";
    });
  }
}

function renderUsersTable(users) {
  const usersTableBody = document.getElementById("usersTableBody");

  if (!users.length) {
    usersTableBody.innerHTML = "<tr><td colspan=\"4\">No users found</td></tr>";
    return;
  }

  usersTableBody.innerHTML = users.map(user => `
    <tr>
      <td>${user.id}</td>
      <td>${user.name}</td>
      <td>${user.username}</td>
      <td>
        <button class="button small" onclick="editUser(${user.id})">Edit</button>
        <button class="button small danger" onclick="deleteUser(${user.id})">Delete</button>
      </td>
    </tr>
  `).join('');
}

function showUserModal(user = null) {
  const modal = document.getElementById("userModal");
  const form = document.getElementById("userForm");
  const modalTitle = document.getElementById("userModalTitle");

  if (user) {
    modalTitle.textContent = "Edit User";
    document.getElementById("userId").value = user.id;
    document.getElementById("userName").value = user.name;
    document.getElementById("userUsername").value = user.username;
    document.getElementById("userPassword").value = ''; // Don't populate password
  } else {
    modalTitle.textContent = "Add User";
    form.reset();
    document.getElementById("userId").value = '';
  }

  modal.style.display = "block";
}

async function editUser(id) {
  try {
    const res = await fetch(`/api/users/${id}`, {
      headers: {
        "Authorization": `Bearer ${getToken()}`
      }
    });
    const user = await res.json();

    if (res.ok) {
      showUserModal(user);
    }
  } catch (error) {
    alert("Failed to load user details");
  }
}

async function deleteUser(id) {
  if (!confirm("Are you sure you want to delete this user?")) return;

  try {
    const res = await fetch(`/api/users/${id}`, {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${getToken()}`
      }
    });

    if (res.ok) {
      loadManageUsersPage(); // Reload the page
    } else {
      alert("Failed to delete user");
    }
  } catch (error) {
    alert("Failed to delete user");
  }
}

async function handleUserSubmit(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const userData = {
    name: formData.get("name"),
    username: formData.get("username"),
    password: formData.get("password")
  };

  const userId = formData.get("userId");
  const isEdit = userId && userId !== '';

  try {
    const res = await fetch(isEdit ? `/api/users/${userId}` : "/api/users", {
      method: isEdit ? "PUT" : "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getToken()}`
      },
      body: JSON.stringify(userData)
    });

    if (res.ok) {
      document.getElementById("userModal").style.display = "none";
      loadManageUsersPage(); // Reload the page
    } else {
      const error = await res.json();
      alert(error.error || "Failed to save user");
    }
  } catch (error) {
    alert("Failed to save user");
  }
}

// Initialize management functionality
const initManagement = () => {
  // Add version button functionality will be added here
};
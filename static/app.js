const tokenKey = "skvms_token";
const userKey = "skvms_user";

const getToken = () => localStorage.getItem(tokenKey);
const setToken = (token) => localStorage.setItem(tokenKey, token);
const clearToken = () => localStorage.removeItem(tokenKey);

const getUser = () => {
  const user = localStorage.getItem(userKey);
  return user ? JSON.parse(user) : null;
};
const setUser = (user) => localStorage.setItem(userKey, JSON.stringify(user));
const clearUser = () => localStorage.removeItem(userKey);

const authBadge = document.getElementById("authBadge");
const logoutBtn = document.getElementById("logoutBtn");
const loginLink = document.getElementById("loginLink");

if (logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    clearToken();
    clearUser();
    window.location.reload();
  });
}

const updateAuthUI = () => {
  const user = getUser();
  if (authBadge) {
    authBadge.textContent = user ? `Welcome, ${user.name}` : "Guest";
  }
  if (loginLink) {
    loginLink.style.display = user ? "none" : "inline-block";
  }
  if (logoutBtn) {
    logoutBtn.style.display = user ? "inline-block" : "none";
  }
};

updateAuthUI();

const loginForm = document.getElementById("loginForm");
if (loginForm) {
  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = new FormData(loginForm);
    const payload = {
      username: formData.get("username"),
      password: formData.get("password"),
    };

    const messageEl = document.getElementById("loginMessage");
    messageEl.textContent = "";
    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      if (!res.ok) {
        messageEl.textContent = data.error || "Login failed";
        return;
      }
      setToken(data.token);
      setUser(data.user);
      updateAuthUI();
      messageEl.textContent = "Login successful. Redirecting...";
      setTimeout(() => {
        window.location.href = "/";
      }, 700);
    } catch (error) {
      messageEl.textContent = "Login failed. Try again.";
    }
  });
}

const deviceList = document.getElementById("deviceList");
if (deviceList) {
  loadDevices();
}

async function loadDevices() {
  deviceList.innerHTML = "<p class=\"muted\">Loading devices...</p>";
  try {
    const res = await fetch("/api/devices");
    const data = await res.json();
    if (!res.ok) {
      deviceList.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
      return;
    }
    renderDevices(data.devices || []);
  } catch (error) {
    deviceList.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
  }
}

function renderDevices(devices) {
  if (!devices.length) {
    deviceList.innerHTML = "<p class=\"muted\">No devices found.</p>";
    return;
  }

  const token = getToken();
  deviceList.innerHTML = devices
    .map(
      (device) => `
      <div class="device-card">
        <div class="device-header">
          <h3>${device.name}</h3>
          <span class="chip">${device.type}</span>
        </div>
        <div class="device-body">
          <p><strong>IP:</strong> ${device.ip_address || "-"}</p>
          <p><strong>MAC:</strong> ${device.mac_address || "-"}</p>
          <p><strong>Firmware:</strong> ${device.firmware_version || "-"}</p>
          <p><strong>Location:</strong> ${device.address || "-"}, ${device.city || "-"}</p>
        </div>
        <div class="device-actions">
          <button class="button" ${token ? "" : "disabled"} data-id="${device.id}">Send Control</button>
          <span class="hint">${token ? "" : "Login required to control"}</span>
        </div>
      </div>
    `
    )
    .join("");

  deviceList.querySelectorAll("button[data-id]").forEach((button) => {
    button.addEventListener("click", async () => {
      const id = button.getAttribute("data-id");
      await sendControlCommand(id);
    });
  });
}

async function sendControlCommand(id) {
  const token = getToken();
  if (!token) {
    alert("Login required to control devices.");
    return;
  }

  try {
    const res = await fetch(`/api/devices/${id}/control`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ command: "connect" }),
    });
    const data = await res.json();
    if (!res.ok) {
      alert(data.error || "Command failed");
      return;
    }
    alert(data.message || "Command sent");
  } catch (error) {
    alert("Command failed");
  }
}
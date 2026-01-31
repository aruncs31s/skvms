// devices.js - Device listing and individual device page functionality

const deviceList = document.getElementById("deviceList");
if (deviceList) {
  loadDevices();
}

const deviceId = getDeviceIdFromPath();
if (deviceId) {
  loadDevicePage(deviceId);
}

const getDeviceIdFromPath = () => {
  const match = window.location.pathname.match(/^\/devices\/(\d+)$/);
  if (!match) return null;
  return match[1];
};

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
      <div class="device-card clickable" data-device-id="${device.id}">
        <div class="device-header">
          <h3>${device.name}</h3>
          <span class="badge">${device.type}</span>
        </div>
        <div class="device-meta">
          <div class="meta-item">
            <strong>IP:</strong> ${device.ip_address || 'N/A'}
          </div>
          <div class="meta-item">
            <strong>MAC:</strong> ${device.mac_address || 'N/A'}
          </div>
          <div class="meta-item">
            <strong>Version:</strong> ${device.firmware_version || 'N/A'}
          </div>
        </div>
        <div class="device-actions">
          <a href="/devices/${device.id}" class="button small">View Details</a>
          ${token ? `<button class="button small secondary control-btn" data-device-id="${device.id}">Control</button>` : ''}
        </div>
      </div>
    `
    )
    .join("");

  // Add click handlers for device cards
  document.querySelectorAll(".device-card.clickable").forEach((card) => {
    card.addEventListener("click", (e) => {
      // Don't navigate if clicking on a button
      if (e.target.tagName === "BUTTON" || e.target.tagName === "A") return;
      const deviceId = card.dataset.deviceId;
      window.location.href = `/devices/${deviceId}`;
    });
  });

  // Add control button handlers
  document.querySelectorAll(".control-btn").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      const deviceId = btn.dataset.deviceId;
      showControlModal(deviceId);
    });
  });
}

async function loadDevicePage(deviceId) {
  try {
    const res = await fetch(`/api/devices/${deviceId}`);
    const device = await res.json();
    if (!res.ok) {
      document.getElementById("deviceTitle").textContent = "Device not found";
      document.getElementById("deviceMeta").innerHTML = "<p class=\"muted\">Failed to load device details.</p>";
      return;
    }

    document.getElementById("deviceTitle").textContent = device.name;
    document.getElementById("deviceMeta").innerHTML = `
      <div class="meta-grid">
        <div class="meta-item"><strong>Type:</strong> ${device.type}</div>
        <div class="meta-item"><strong>IP Address:</strong> ${device.ip_address || 'N/A'}</div>
        <div class="meta-item"><strong>MAC Address:</strong> ${device.mac_address || 'N/A'}</div>
        <div class="meta-item"><strong>Firmware:</strong> ${device.firmware_version || 'N/A'}</div>
        <div class="meta-item"><strong>Address:</strong> ${device.address || 'N/A'}</div>
        <div class="meta-item"><strong>City:</strong> ${device.city || 'N/A'}</div>
      </div>
    `;

    // Load readings for this device
    loadDeviceReadings(deviceId);
  } catch (error) {
    document.getElementById("deviceTitle").textContent = "Error";
    document.getElementById("deviceMeta").innerHTML = "<p class=\"muted\">Failed to load device details.</p>";
  }
}

async function loadDeviceReadings(deviceId) {
  const readingsContainer = document.getElementById("readingsContainer");
  if (!readingsContainer) return;

  readingsContainer.innerHTML = "<p class=\"muted\">Loading readings...</p>";

  try {
    const datePicker = document.getElementById("datePicker");
    const showYesterday = document.getElementById("showYesterday");
    const historyLimit = document.getElementById("historyLimit");

    let url = `/api/devices/${deviceId}/readings?limit=${historyLimit.value}`;

    if (datePicker && datePicker.value) {
      url += `&date=${datePicker.value}`;
    }

    if (showYesterday && showYesterday.checked) {
      url += `&yesterday=true`;
    }

    const res = await fetch(url);
    const data = await res.json();

    if (!res.ok) {
      readingsContainer.innerHTML = "<p class=\"muted\">Failed to load readings.</p>";
      return;
    }

    renderReadings(data.readings || []);
  } catch (error) {
    readingsContainer.innerHTML = "<p class=\"muted\">Failed to load readings.</p>";
  }
}

function renderReadings(readings) {
  const readingsContainer = document.getElementById("readingsContainer");
  if (!readings.length) {
    readingsContainer.innerHTML = "<p class=\"muted\">No readings found.</p>";
    return;
  }

  // Create chart
  const chartData = readings.map(r => ({
    timestamp: new Date(r.timestamp * 1000),
    voltage: r.voltage,
    current: r.current
  })).reverse();

  Highcharts.chart('readingsChart', {
    title: { text: 'Voltage & Current Readings' },
    xAxis: { type: 'datetime' },
    yAxis: [{
      title: { text: 'Voltage (V)' },
      opposite: false
    }, {
      title: { text: 'Current (A)' },
      opposite: true
    }],
    series: [{
      name: 'Voltage',
      data: chartData.map(d => [d.timestamp.getTime(), d.voltage]),
      yAxis: 0,
      color: '#007bff'
    }, {
      name: 'Current',
      data: chartData.map(d => [d.timestamp.getTime(), d.current]),
      yAxis: 1,
      color: '#28a745'
    }]
  });

  // Create readings table
  readingsContainer.innerHTML = `
    <div id="readingsChart" style="height: 400px; margin-bottom: 20px;"></div>
    <table class="table">
      <thead>
        <tr>
          <th>Time</th>
          <th>Voltage (V)</th>
          <th>Current (A)</th>
          <th>Power (W)</th>
        </tr>
      </thead>
      <tbody>
        ${readings.map(reading => `
          <tr>
            <td>${new Date(reading.timestamp * 1000).toLocaleString()}</td>
            <td>${reading.voltage.toFixed(2)}</td>
            <td>${reading.current.toFixed(2)}</td>
            <td>${(reading.voltage * reading.current).toFixed(2)}</td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
}

// Initialize device-related event listeners
const initDevices = () => {
  // Refresh button
  const refreshBtn = document.getElementById("refreshBtn");
  if (refreshBtn) {
    refreshBtn.addEventListener("click", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        loadDeviceReadings(deviceId);
      }
    });
  }

  // Date picker
  const datePicker = document.getElementById("datePicker");
  if (datePicker) {
    datePicker.addEventListener("change", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        loadDeviceReadings(deviceId);
      }
    });
  }

  // Show yesterday checkbox
  const showYesterday = document.getElementById("showYesterday");
  if (showYesterday) {
    showYesterday.addEventListener("change", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        loadDeviceReadings(deviceId);
      }
    });
  }

  // History limit selector
  const historyLimit = document.getElementById("historyLimit");
  if (historyLimit) {
    historyLimit.addEventListener("change", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        loadDeviceReadings(deviceId);
      }
    });
  }

  // Export buttons
  const exportCsvBtn = document.getElementById("exportCsvBtn");
  if (exportCsvBtn) {
    exportCsvBtn.addEventListener("click", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        exportReadingsCSV(deviceId);
      }
    });
  }

  const exportJsonBtn = document.getElementById("exportJsonBtn");
  if (exportJsonBtn) {
    exportJsonBtn.addEventListener("click", () => {
      const deviceId = getDeviceIdFromPath();
      if (deviceId) {
        exportReadingsJSON(deviceId);
      }
    });
  }

  // Version modal functionality
  const addVersionBtn = document.getElementById("addVersionBtn");
  const versionModal = document.getElementById("versionModal");
  const closeVersionModal = document.getElementById("closeVersionModal");
  const cancelVersionBtn = document.getElementById("cancelVersionBtn");
  const versionForm = document.getElementById("versionForm");
  const addFeatureBtn = document.getElementById("addFeatureBtn");
  const featuresContainer = document.getElementById("featuresContainer");

  if (addVersionBtn) {
    addVersionBtn.addEventListener("click", () => {
      const token = getToken();
      if (!token) {
        alert("Please login first");
        return;
      }
      openVersionModal();
    });
  }

  if (closeVersionModal) {
    closeVersionModal.addEventListener("click", () => {
      versionModal.style.display = "none";
    });
  }

  if (cancelVersionBtn) {
    cancelVersionBtn.addEventListener("click", () => {
      versionModal.style.display = "none";
    });
  }

  if (addFeatureBtn) {
    addFeatureBtn.addEventListener("click", () => {
      addFeatureField();
    });
  }

  if (versionForm) {
    versionForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      await createVersion();
    });
  }

  // Close modal when clicking outside
  if (versionModal) {
    versionModal.addEventListener("click", (e) => {
      if (e.target === versionModal) {
        versionModal.style.display = "none";
      }
    });
  }

  // Remove feature functionality
  if (featuresContainer) {
    featuresContainer.addEventListener("click", (e) => {
      if (e.target.classList.contains("remove-feature")) {
        e.target.closest(".feature-item").remove();
      }
    });
  }
};

// Version modal functions
function openVersionModal() {
  const versionModal = document.getElementById("versionModal");
  const versionForm = document.getElementById("versionForm");
  const featuresContainer = document.getElementById("featuresContainer");

  // Reset form
  versionForm.reset();

  // Clear existing features except the first one
  const featureItems = featuresContainer.querySelectorAll(".feature-item");
  for (let i = 1; i < featureItems.length; i++) {
    featureItems[i].remove();
  }

  // Reset first feature
  const firstFeature = featuresContainer.querySelector(".feature-item");
  if (firstFeature) {
    firstFeature.querySelector(".feature-name").value = "";
    firstFeature.querySelector(".feature-enabled").checked = true;
  }

  versionModal.style.display = "flex";
}

function addFeatureField() {
  const featuresContainer = document.getElementById("featuresContainer");
  const featureItem = document.createElement("div");
  featureItem.className = "feature-item";
  featureItem.innerHTML = `
    <input type="text" class="feature-name" placeholder="Feature name" required>
    <label class="checkbox-label">
      <input type="checkbox" class="feature-enabled" checked>
      Enabled
    </label>
    <button type="button" class="button ghost remove-feature">Remove</button>
  `;
  featuresContainer.appendChild(featureItem);
}

async function createVersion() {
  const token = getToken();
  if (!token) {
    alert("Please login first");
    return;
  }

  const versionNumber = document.getElementById("versionNumber").value;
  const featureItems = document.querySelectorAll(".feature-item");

  const features = [];
  featureItems.forEach(item => {
    const name = item.querySelector(".feature-name").value.trim();
    const enabled = item.querySelector(".feature-enabled").checked;
    if (name) {
      features.push({ name, enabled });
    }
  });

  if (!versionNumber || features.length === 0) {
    alert("Please fill in version number and at least one feature");
    return;
  }

  try {
    // First create the version
    const versionRes = await fetch("/api/versions", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        version: versionNumber,
      }),
    });

    const versionData = await versionRes.json();
    if (!versionRes.ok) {
      alert(versionData.error || "Failed to create version");
      return;
    }

    // Then create the features
    for (const feature of features) {
      const featureRes = await fetch(`/api/features`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          version_id: versionData.id,
          name: feature.name,
          enabled: feature.enabled
        }),
      });

      if (!featureRes.ok) {
        const featureData = await featureRes.json();
        alert(`Failed to create feature "${feature.name}": ${featureData.error || "Unknown error"}`);
        return;
      }
    }

    alert("Version created successfully!");
    document.getElementById("versionModal").style.display = "none";
  } catch (error) {
    alert("Failed to create version");
  }
}
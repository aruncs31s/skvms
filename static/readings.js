// readings.js - All readings page functionality

const allReadingsPage = window.location.pathname === "/all-readings";
if (allReadingsPage) {
  loadAllReadingsPage();
}

async function loadAllReadingsPage() {
  const readingsContainer = document.getElementById("allReadingsContainer");
  if (!readingsContainer) return;

  readingsContainer.innerHTML = "<p class=\"muted\">Loading all readings...</p>";

  try {
    const res = await fetch("/api/devices");
    const data = await res.json();

    if (!res.ok) {
      readingsContainer.innerHTML = "<p class=\"muted\">Failed to load devices.</p>";
      return;
    }

    const devices = data.devices || [];
    if (!devices.length) {
      readingsContainer.innerHTML = "<p class=\"muted\">No devices found.</p>";
      return;
    }

    // Load readings for all devices
    const allReadings = [];
    for (const device of devices) {
      try {
        const readingsRes = await fetch(`/api/devices/${device.id}/readings?limit=10`);
        const readingsData = await readingsRes.json();
        if (readingsRes.ok && readingsData.readings) {
          readingsData.readings.forEach(reading => {
            allReadings.push({
              ...reading,
              deviceName: device.name,
              deviceId: device.id
            });
          });
        }
      } catch (error) {
        console.error(`Failed to load readings for device ${device.id}:`, error);
      }
    }

    // Sort by timestamp descending
    allReadings.sort((a, b) => b.timestamp - a.timestamp);

    renderAllReadings(allReadings);
  } catch (error) {
    readingsContainer.innerHTML = "<p class=\"muted\">Failed to load readings.</p>";
  }
}

function renderAllReadings(readings) {
  const readingsContainer = document.getElementById("allReadingsContainer");

  if (!readings.length) {
    readingsContainer.innerHTML = "<p class=\"muted\">No readings found.</p>";
    return;
  }

  readingsContainer.innerHTML = `
    <table class="table">
      <thead>
        <tr>
          <th>Device</th>
          <th>Time</th>
          <th>Voltage (V)</th>
          <th>Current (A)</th>
          <th>Power (W)</th>
        </tr>
      </thead>
      <tbody>
        ${readings.map(reading => `
          <tr>
            <td><a href="/devices/${reading.deviceId}">${reading.deviceName}</a></td>
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
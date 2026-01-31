// controls.js - Device control functionality

async function sendControlCommand(id) {
  const command = prompt("Enter control command (e.g., 'restart', 'shutdown', 'ledon', 'ledoff'):");
  if (!command) return;

  try {
    const res = await fetch(`/api/devices/${id}/control`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${getToken()}`
      },
      body: JSON.stringify({ command })
    });

    const data = await res.json();
    if (res.ok) {
      alert(`Command "${command}" sent successfully`);
    } else {
      alert(data.error || "Failed to send command");
    }
  } catch (error) {
    alert("Failed to send command");
  }
}

function showControlModal(deviceId) {
  const command = prompt("Enter control command (e.g., 'restart', 'shutdown', 'ledon', 'ledoff'):");
  if (!command) return;

  sendControlCommand(deviceId);
}

// Export functionality for readings
async function exportReadingsCSV(deviceId, startTime, endTime, limit) {
  try {
    let url = `/api/devices/${deviceId}/readings?format=csv&limit=${limit || 1000}`;
    if (startTime) url += `&start=${startTime}`;
    if (endTime) url += `&end=${endTime}`;

    const res = await fetch(url);
    if (!res.ok) throw new Error('Export failed');

    const blob = await res.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = `device-${deviceId}-readings-${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    window.URL.revokeObjectURL(downloadUrl);
  } catch (error) {
    alert("Failed to export CSV");
  }
}

async function exportReadingsJSON(deviceId, startTime, endTime, limit) {
  try {
    let url = `/api/devices/${deviceId}/readings?format=json&limit=${limit || 1000}`;
    if (startTime) url += `&start=${startTime}`;
    if (endTime) url += `&end=${endTime}`;

    const res = await fetch(url);
    if (!res.ok) throw new Error('Export failed');

    const blob = await res.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    link.href = downloadUrl;
    link.download = `device-${deviceId}-readings-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
    window.URL.revokeObjectURL(downloadUrl);
  } catch (error) {
    alert("Failed to export JSON");
  }
}
// app.js - Main application entry point

// Include all module scripts
document.addEventListener('DOMContentLoaded', function() {
  // Load modules in order
  const modules = [
    'auth.js',
    'navbar.js',
    'devices.js',
    'management.js',
    'readings.js',
    'audit.js',
    'controls.js'
  ];

  let loadedCount = 0;

  modules.forEach(module => {
    const script = document.createElement('script');
    script.src = `/static/${module}`;
    script.onload = function() {
      loadedCount++;
      if (loadedCount === modules.length) {
        // All modules loaded, initialize
        initializeApp();
      }
    };
    script.onerror = function() {
      console.error(`Failed to load ${module}`);
    };
    document.head.appendChild(script);
  });
});

function initializeApp() {
  // Initialize auth UI
  updateAuthUI();

  // Initialize page-specific functionality
  if (typeof initAuth === 'function') initAuth();
  if (typeof initDevices === 'function') initDevices();
  if (typeof initManagement === 'function') initManagement();
}
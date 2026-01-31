// auth.js - Authentication utilities and login functionality

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

// Update authentication UI elements
function updateAuthUI() {
  const user = getUser();
  const displayText = user ? `Welcome, ${user.name}` : "Guest";

  // Update navbar elements
  const navAuthBadge = document.getElementById("navAuthBadge");
  const navLoginLink = document.getElementById("navLoginLink");
  const navLogoutBtn = document.getElementById("navLogoutBtn");
  const navManageDevices = document.getElementById("navManageDevices");
  const navManageUsers = document.getElementById("navManageUsers");
  const navAudit = document.getElementById("navAudit");

  if (navAuthBadge) navAuthBadge.textContent = displayText;
  if (navLoginLink) navLoginLink.style.display = user ? "none" : "inline-block";
  if (navLogoutBtn) navLogoutBtn.style.display = user ? "inline-block" : "none";

  // Show admin links only when logged in
  if (navManageDevices) navManageDevices.style.display = user ? "inline-block" : "none";
  if (navManageUsers) navManageUsers.style.display = user ? "inline-block" : "none";
  if (navAudit) navAudit.style.display = user ? "inline-block" : "none";

  // Update other auth elements
  const authBadge = document.getElementById("authBadge");
  const loginLink = document.getElementById("loginLink");
  const logoutBtn = document.getElementById("logoutBtn");

  if (authBadge) authBadge.textContent = displayText;
  if (loginLink) loginLink.style.display = user ? "none" : "inline-block";
  if (logoutBtn) logoutBtn.style.display = user ? "inline-block" : "none";

  // Mark active page in navbar
  const path = window.location.pathname;
  document.querySelectorAll(".navbar-item").forEach((item) => {
    const href = item.getAttribute("href");
    if (href && (href === path || (href !== "/" && path.startsWith(href)))) {
      item.classList.add("active");
    }
  });
}

// Login form handling
const initLoginForm = () => {
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
};

// Initialize auth-related elements
const initAuth = () => {
  initLoginForm();
};
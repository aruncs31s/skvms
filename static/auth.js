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
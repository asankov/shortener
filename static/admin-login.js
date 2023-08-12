document.getElementById("login-form").onsubmit = async e => {
  e.preventDefault();

  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;

  const response = await fetch("/admin/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  const json = await response.json();
  localStorage.setItem("token", json.token);

  window.location.replace("/admin")
};

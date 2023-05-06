urlTypeButton = document.getElementById("url-type-button");

document.getElementById("random-url-button").onclick = () => {
  urlTypeButton.innerHTML = "Random URL";
};
document.getElementById("custom-url-button").onclick = () => {
  urlTypeButton.innerHTML = "Custom URL";
};

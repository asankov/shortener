urlTypeButton = document.getElementById("url-type-button");

document.getElementById("random-url-button").onclick = () => {
  urlTypeButton.innerHTML = "Random URL";
};
document.getElementById("custom-url-button").onclick = () => {
  urlTypeButton.innerHTML = "Custom URL";
};
document.getElementById("shorten-link-form").onsubmit = e => {
  e.preventDefault();

  const link = document.getElementById("link").value;

  console.log("submitted");
  const response = fetch("/api/v1/links", {
    method: "POST",
    body: JSON.stringify({ url: link }),
  });

  response
    .then(() => {
      // TODO: success banner
      location.reload();
    })
    .catch(e => {
      // TODO: proper error msg
      console.log("ERROR", e);
    });
};

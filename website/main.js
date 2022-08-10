const localFuncApi = "http://localhost:7071/api/GetVisitorCounter"

function visitorCount() {
  let counter = document.getElementById('counter').innerText;
  let count = 0;
  count = parseInt(counter.innerText);
  counter.innerHTML = count;
}

window.onload = visitorCount();

async function main() {
  const response = await fetch("https://resume.acaldwell.dev/resume.json", {
    method: "GET"
  });

  if (response.ok) {
    let json = await response.json();
  } else {
    console.log("HTTP-Error: " + response.status);
  }
}

main();
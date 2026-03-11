const div = document.getElementById("coin-matrix");
const rect = div.getBoundingClientRect();

const width = rect.width;
const height = rect.height;

const img = div.querySelector(".img-container");
const fresh = img.cloneNode(true);

const furu = (ele) => {
  // time ms
  const minTime = 5000;
  const maxTime = 10000;
  let elapsed = 0;
  let intervalId;
  const duration =
    Math.floor(Math.random() * (maxTime - minTime + 1)) + minTime;

  intervalId = setInterval(() => {
    elapsed += 10;
    ele.style.bottom = `${((duration - elapsed) / duration) * 100}%`;

    if (elapsed >= duration + 5000) {
      clearInterval(intervalId);
      ele.remove();
    }

    ele.onclick = () => {
      clearInterval(intervalId);
      ele.remove();
      newCoin(0);
    };
  }, 10);
};

const newCoin = (ms, recursive) => {
  setTimeout(() => {
    const nc = fresh.cloneNode(true);

    // size
    const minW = 5;
    const maxW = 15;
    const factorW =
      (Math.floor(Math.random() * (maxW - minW + 1)) + minW) / 100;
    const size = ((width * factorW) / width) * 100;
    nc.style.width = `${size}%`;

    // x
    const minX = 0;
    const maxX = 100 - size;
    const factorX = Math.floor(Math.random() * (maxX - minX + 1)) + minX;
    const x = (width * factorX) / width;

    nc.style.left = `${x}%`;
    nc.style.transform = "translateY(-100%)";
    furu(nc);
    div.appendChild(nc);
    if (recursive) {
      newCoin(700, true);
    }
  }, ms);
};
newCoin(0);
newCoin(100);
newCoin(200);
newCoin(300);
newCoin(400);
newCoin(0, true);

const R1 = 1;
const R2 = 2;
const K1 = 150;
const K2 = 5;
let zBuffer = [];
let output = [];
let A = 0;
let B = 0;

function setup() {
  createCanvas(400, 400);

  for (let y = 0; y < height; y++) {
    let bufferRow = [];
    let outputRow = [];
    for (let x = 0; x < width; x++) {
      bufferRow.push(0);
      outputRow.push(" ");
    }
    zBuffer.push(bufferRow);
    output.push(outputRow);
  }
}

function draw() {
  background(0);
  noStroke();
  renderFrame(A, B);
  A += 0.04;
  B += 0.02;
}

function renderFrame(A, B) {
  resetArrays();

  let [cosA, sinA, cosB, sinB] = [
    Math.cos(A),
    Math.sin(A),
    Math.cos(B),
    Math.sin(B),
  ];

  for (let theta = 0; theta < 6.28; theta += 0.3) {
    let [cosTheta, sinTheta] = [Math.cos(theta), Math.sin(theta)];
    for (let phi = 0; phi < 6.28; phi += 0.1) {
      let [cosPhi, sinPhi] = [Math.cos(phi), Math.sin(phi)];

      let x =
        (R2 + R1 * cosTheta) * (cosB * cosPhi + sinA * sinB * sinPhi) -
        R1 * cosA * sinB * sinTheta;
      let y =
        (R2 + R1 * cosTheta) * (sinB * cosPhi - cosB * sinA * sinPhi) +
        R1 * cosA * cosB * sinTheta;
      let z = cosA * (R2 + R1 * cosTheta) * sinPhi + R1 * sinA * sinTheta + K2;

      let ooz = 1 / z;

      let xp = parseInt(width / 2 + K1 * x * ooz);
      let yp = parseInt(height / 2 - K1 * y * ooz);

      if (xp < 0 || xp >= width || yp < 0 || yp >= height) {
        continue;
      }

      // Calculated from L = (Nx, Ny, Nz) dot (0, 1, -1) <- pre-chosen light vector
      let luminance =
        cosPhi * cosTheta * sinB -
        cosA * cosTheta * sinPhi -
        sinA * sinTheta +
        cosB * (cosA * sinTheta - cosTheta * sinA * sinPhi);
      if (luminance <= 0) {
        continue;
      }

      // larger 1/z means pixel is closer so it should override for current x', y'
      if (ooz <= zBuffer[yp][xp]) {
        continue;
      }

      let luminance_index = int(luminance * 8);
      if (luminance_index < 0) {
        continue;
      }

      zBuffer[yp][xp] = ooz;
      // output[yp][xp] = string(".,-~:;=!*#$@"[luminance_index]);
      circle(xp, yp, 1.5);
      fill(255, 255, 255, luminance_index * 30);
    }
  }
}

function resetArrays() {
  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      zBuffer[y][x] = 0;
      output[y][x] = " ";
    }
  }
}

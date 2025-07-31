const canvas = document.getElementById('drawCanvas');
const ctx = canvas.getContext('2d');

canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

let isDrawing = false;
let lastX = 0;
let lastY = 0;
let selectedShape = 'free';
let selectedColor = 'black';
let shapes = [];

// Shape options
const rectButton = document.getElementById('rect');
const circleButton = document.getElementById('circle');
const lineButton = document.getElementById('line');
const freeButton = document.getElementById('free');
const eraseButton = document.getElementById('erase');

rectButton.addEventListener('click', () => selectedShape = 'rect');
circleButton.addEventListener('click', () => selectedShape = 'circle');
lineButton.addEventListener('click', () => selectedShape = 'line');
freeButton.addEventListener('click', () => selectedShape = 'free');
eraseButton.addEventListener('click', () => selectedShape = 'erase');

// Color options
const redButton = document.getElementById('red');
const greenButton = document.getElementById('green');
const blueButton = document.getElementById('blue');

redButton.addEventListener('click', () => selectedColor = 'red');
greenButton.addEventListener('click', () => selectedColor = 'green');
blueButton.addEventListener('click', () => selectedColor = 'blue');

function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    shapes.forEach(shape => {
        ctx.beginPath();
        ctx.strokeStyle = shape.color;
        if (shape.type === 'line') {
            ctx.moveTo(shape.startX, shape.startY);
            ctx.lineTo(shape.endX, shape.endY);
            ctx.stroke();
        } else if (shape.type === 'rect') {
            ctx.rect(shape.startX, shape.startY, shape.width, shape.height);
            ctx.stroke();
        } else if (shape.type === 'circle') {
            ctx.arc(shape.startX, shape.startY, shape.radius, 0, 2 * Math.PI);
            ctx.stroke();
        }
    });
}

function handleMouseDown(e) {
    isDrawing = true;
    lastX = e.clientX;
    lastY = e.clientY;
}

function handleMouseUp(e) {
    isDrawing = false;
    if (selectedShape === 'free') return;
    const shape = { color: selectedColor, startX: lastX, startY: lastY };
    if (selectedShape === 'line') {
        shape.type = 'line';
        shape.endX = e.clientX;
        shape.endY = e.clientY;
    } else if (selectedShape === 'rect') {
        shape.type = 'rect';
        shape.width = e.clientX - lastX;
        shape.height = e.clientY - lastY;
    } else if (selectedShape === 'circle') {
        shape.type = 'circle';
        shape.radius = Math.sqrt(Math.pow(e.clientX - lastX, 2) + Math.pow(e.clientY - lastY, 2));
    }
    shapes.push(shape);
    draw();
}

function handleMouseMove(e) {
    if (!isDrawing) return;
    if (selectedShape === 'free') {
        const shape = { color: selectedColor, startX: lastX, startY: lastY, endX: e.clientX, endY: e.clientY, type: 'line' };
        shapes.push(shape);
        draw();
        lastX = e.clientX;
        lastY = e.clientY;
    }
}

canvas.addEventListener('mousedown', handleMouseDown);
canvas.addEventListener('mouseup', handleMouseUp);
canvas.addEventListener('mousemove', handleMouseMove);

// script.js
// Main whiteboard class
class Whiteboard {
    constructor() {
        this.canvas = document.getElementById('whiteboard');
        this.ctx = this.canvas.getContext('2d');
        this.shapePreview = document.getElementById('shape-preview');
        
        this.currentTool = 'draw';
        this.currentColor = '#000000';
        this.currentThickness = 3;
        this.eraserSize = 20;
        this.isDrawing = false;
        this.isResizing = false;
        this.startX = 0;
        this.startY = 0;
        this.currentObject = null;
        this.selectedObject = null;
        this.objects = [];
        this.resizeHandleIndex = -1;
        
        // WebSocket connection
        this.ws = null;
        this.connectWebSocket();
        
        this.init();
    }

    async connectWebSocket() {
        // Create WebSocket connection
        const token = await this.fetchToken()
        if(!token){
            window.location.href="/login/"
        }

        const host = window.location.host
        this.ws = new WebSocket(`ws://${host}/ws?token=${token}`);
        
        // Add WebSocket event listeners
        this.ws.onopen = () => console.log('WebSocket connected');
        this.ws.onmessage = (event) => this.handleWebSocketMessage(event);
        this.ws.onclose = () => console.log('WebSocket disconnected');
        this.ws.onerror = (error) => console.error('WebSocket error:', error);
    }

    handleWebSocketMessage(event) {
        const message = JSON.parse(event.data);
        if(message.type == "TOKEN_EXPIRED"){
            this.reconnect();
        }else{
            // Process incoming drawing events from other users
            this.processRemoteEvent(message);
        }
    }

    async fetchToken(){
        const token = localStorage.getItem("auth-token")
        if(token){
            return token
        }

        const refreshed = await this.refreshAccessToken();
        if (refreshed) {
            window.location.href = "/drawboard";
        }
    }

    async refreshAccessToken(){
        try{
            const res =  await fetch('/refresh', {
                method: 'POST',
                credentials: 'include' // send the refresh token cookie
            });

            if (!res.ok) return false;

            const data = await res.json();
            localStorage.setItem('auth-token', data["auth-token"]);
            localStorage.setItem('token-expiry', Date.now() + 15 * 60 * 1000);
            return true;
        }catch(error){
            console.error("Refresh token request failed:", error);
            return false;
        }
    }

    reconnect() {
        this.ws.close();
        this.ws = null;
        this.connectWS();
    }

    sendDrawingEvent(eventData) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(eventData));
        }
    }

    processRemoteEvent(eventData) {
        // Process events from other users
        // This would redraw the canvas based on remote events
        console.log('Processing remote event:', eventData);
        
        switch (eventData.type) {
            case 'freehandDraw':
                // Create a new object based on the event data
                const drawObj = {
                    type: eventData.tool,
                    color: eventData.data.color,
                    thickness: eventData.data.thickness,
                    points: eventData.data.points
                };
                this.objects.push(drawObj);
                this.redraw();
                break;
                
            case 'shapeCreate':
                // Create a new shape object based on the event data
                const shapeObj = {
                    type: eventData.tool,
                    color: eventData.data.color,
                    thickness: eventData.data.thickness,
                    x: eventData.data.x,
                    y: eventData.data.y,
                    width: eventData.data.width,
                    height: eventData.data.height
                };
                this.objects.push(shapeObj);
                this.redraw();
                break;
                
            case 'textAdd':
                // Create a new text object based on the event data
                const textObj = {
                    type: 'text',
                    color: eventData.data.color,
                    thickness: eventData.data.thickness,
                    x: eventData.data.x,
                    y: eventData.data.y,
                    width: 0,
                    height: 0,
                    text: eventData.data.text
                };
                this.objects.push(textObj);
                this.redraw();
                break;
                
            case 'objectDelete':
                // Remove object at the specified index
                if (eventData.data.index < this.objects.length) {
                    this.objects.splice(eventData.data.index, 1);
                    this.redraw();
                }
                break;
                
            case 'boardClear':
                // Clear the entire board
                this.objects = [];
                this.deselectObject();
                this.redraw();
                break;
                
            default:
                console.log('Unknown event type:', eventData.type);
        }
    }

    init() {
        this.resizeCanvas();
        window.addEventListener('resize', () => this.resizeCanvas());
        
        this.canvas.addEventListener('mousedown', (e) => this.handleMouseDown(e));
        this.canvas.addEventListener('mousemove', (e) => this.handleMouseMove(e));
        this.canvas.addEventListener('mouseup', () => this.handleMouseUp());
        this.canvas.addEventListener('mouseout', () => this.handleMouseUp());
        
        document.querySelectorAll('.tool').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.setTool(e.target.dataset.tool);
            });
        });
        
        document.querySelector('.color-picker').addEventListener('input', (e) => {
            this.currentColor = e.target.value;
        });
        
        document.getElementById('thickness').addEventListener('input', (e) => {
            this.currentThickness = e.target.value;
            document.getElementById('thickness-value').textContent = `${e.target.value}px`;
        });
        
        document.getElementById('eraser-size').addEventListener('input', (e) => {
            this.eraserSize = e.target.value;
            document.getElementById('eraser-size-value').textContent = `${e.target.value}px`;
        });
        
        document.getElementById('clear-board').addEventListener('click', () => {
            if (confirm('Are you sure you want to clear the entire board?')) {
                this.clearBoard();
            }
        });

        // Add text input handling
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Delete' || e.key === 'Backspace') {
                if (this.selectedObject && this.currentTool === 'select') {
                    this.deleteSelectedObject();
                }
            }
        });
    }

    resizeCanvas() {
        this.canvas.width = this.canvas.offsetWidth;
        this.canvas.height = this.canvas.offsetHeight;
        this.redraw();
    }

    setTool(tool) {
        this.currentTool = tool;
        document.querySelectorAll('.tool').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-tool="${tool}"]`).classList.add('active');
        
        // Update cursor
        this.updateCursor();
        
        // Show/hide eraser size controls
        const eraserSizeLabel = document.getElementById('eraser-size-label');
        const eraserSizeSlider = document.getElementById('eraser-size');
        const eraserSizeValue = document.getElementById('eraser-size-value');
        
        if (tool === 'erase') {
            eraserSizeLabel.style.display = 'inline';
            eraserSizeSlider.style.display = 'inline';
            eraserSizeValue.style.display = 'inline';
            // Hide pen size controls
            document.getElementById('thickness').style.display = 'none';
            document.getElementById('thickness-value').style.display = 'none';
            document.querySelector('label[for="thickness"]').style.display = 'none';
        } else {
            eraserSizeLabel.style.display = 'none';
            eraserSizeSlider.style.display = 'none';
            eraserSizeValue.style.display = 'none';
            // Show pen size controls
            document.getElementById('thickness').style.display = 'inline';
            document.getElementById('thickness-value').style.display = 'inline';
            document.querySelector('label[for="thickness"]').style.display = 'inline';
        }
        
        if (tool !== 'select') {
            this.deselectObject();
        }
    }

    updateCursor() {
        // Remove all tool cursor classes
        this.canvas.classList.remove(
            'draw-cursor', 
            'erase-cursor', 
            'line-cursor', 
            'rectangle-cursor', 
            'circle-cursor', 
            'text-cursor', 
            'select-cursor'
        );
        
        // Add cursor class based on current tool
        this.canvas.classList.add(`${this.currentTool}-cursor`);
    }

    handleMouseDown(e) {
        console.log('Mouse down event:', {
            tool: this.currentTool,
            x: this.startX,
            y: this.startY
        });
        
        const rect = this.canvas.getBoundingClientRect();
        this.startX = e.clientX - rect.left;
        this.startY = e.clientY - rect.top;
        
        if (this.currentTool === 'select') {
            this.handleSelectMouseDown();
        } else if (this.currentTool === 'text') {
            this.addTextObject(this.startX, this.startY);
        } else {
            this.isDrawing = true;
            this.currentObject = {
                type: this.currentTool,
                x: this.startX,
                y: this.startY,
                width: 0,
                height: 0,
                color: this.currentTool === 'erase' ? '#FFFFFF' : this.currentColor,
                thickness: this.currentTool === 'erase' ? this.eraserSize : this.currentThickness,
                points: this.currentTool === 'draw' || this.currentTool === 'erase' ? [{ x: this.startX, y: this.startY }] : []
            };
        }
    }

    handleMouseMove(e) {
        if (!this.isDrawing && !this.isResizing && this.currentTool !== 'select') return;
        
        const rect = this.canvas.getBoundingClientRect();
        const currentX = e.clientX - rect.left;
        const currentY = e.clientY - rect.top;
        
        if (this.isResizing && this.selectedObject) {
            this.resizeObject(currentX, currentY);
            return;
        }
        
        if (this.isDrawing) {
            if (this.currentTool === 'draw' || this.currentTool === 'erase') {
                this.currentObject.points.push({ x: currentX, y: currentY });
                this.redraw();
                this.drawObject(this.currentObject);
            } else {
                this.currentObject.width = currentX - this.startX;
                this.currentObject.height = currentY - this.startY;
                this.redraw();
                this.drawObject(this.currentObject);
                
                // Show preview for shapes
                this.updateShapePreview(this.startX, this.startY, currentX, currentY);
            }
        } else if (this.currentTool === 'select') {
            this.handleSelectMouseMove(currentX, currentY);
        }
    }

    handleMouseUp() {
        console.log('Mouse up event triggered');
        
        if (!this.isDrawing && !this.isResizing) {
            console.log('No drawing or resizing in progress');
            return;
        }
        
        if (this.isDrawing) {
            if ((this.currentTool === 'draw' || this.currentTool === 'erase') && this.currentObject.points.length > 1) {
                console.log('Freehand drawing completed:', {
                    tool: this.currentTool,
                    pointsCount: this.currentObject.points.length,
                    color: this.currentObject.color,
                    thickness: this.currentObject.thickness
                });
                this.objects.push(this.currentObject);
                
                // Send drawing event through WebSocket
                this.sendDrawingEvent({
                    type: 'freehandDraw',
                    tool: this.currentTool,
                    data: {
                        color: this.currentObject.color,
                        thickness: this.currentObject.thickness,
                        points: this.currentObject.points
                    }
                });
            } else if (this.currentTool !== 'draw' && this.currentTool !== 'erase') {
                // Adjust negative dimensions
                if (this.currentObject.width < 0) {
                    this.currentObject.x += this.currentObject.width;
                    this.currentObject.width = Math.abs(this.currentObject.width);
                }
                if (this.currentObject.height < 0) {
                    this.currentObject.y += this.currentObject.height;
                    this.currentObject.height = Math.abs(this.currentObject.height);
                }
                
                // Only add if shape has meaningful size
                if (Math.abs(this.currentObject.width) > 5 || Math.abs(this.currentObject.height) > 5) {
                    console.log('Shape created:', {
                        type: this.currentTool,
                        x: this.currentObject.x,
                        y: this.currentObject.y,
                        width: this.currentObject.width,
                        height: this.currentObject.height,
                        color: this.currentObject.color,
                        thickness: this.currentObject.thickness
                    });
                    this.objects.push(this.currentObject);
                    
                    // Send shape creation event through WebSocket
                    this.sendDrawingEvent({
                        type: 'shapeCreate',
                        tool: this.currentTool,
                        data: {
                            color: this.currentObject.color,
                            thickness: this.currentObject.thickness,
                            x: this.currentObject.x,
                            y: this.currentObject.y,
                            width: this.currentObject.width,
                            height: this.currentObject.height
                        }
                    });
                }
            }
            this.currentObject = null;
        }
        
        this.isDrawing = false;
        this.isResizing = false;
        this.shapePreview.style.display = 'none';
        this.redraw();
    }

    updateShapePreview(startX, startY, currentX, currentY) {
        const width = currentX - startX;
        const height = currentY - startY;
        
        this.shapePreview.style.display = 'block';
        this.shapePreview.style.left = `${Math.min(startX, currentX)}px`;
        this.shapePreview.style.top = `${Math.min(startY, currentY)}px`;
        this.shapePreview.style.width = `${Math.abs(width)}px`;
        this.shapePreview.style.height = `${Math.abs(height)}px`;
    }

    drawObject(obj) {
        this.ctx.strokeStyle = obj.color;
        this.ctx.fillStyle = obj.color;
        this.ctx.lineWidth = obj.thickness;
        this.ctx.lineCap = 'round';
        this.ctx.lineJoin = 'round';
        
        switch (obj.type) {
            case 'draw':
            case 'erase':
                if (obj.points.length < 2) return;
                this.ctx.beginPath();
                this.ctx.moveTo(obj.points[0].x, obj.points[0].y);
                for (let i = 1; i < obj.points.length; i++) {
                    this.ctx.lineTo(obj.points[i].x, obj.points[i].y);
                }
                this.ctx.stroke();
                break;
                
            case 'line':
                this.ctx.beginPath();
                this.ctx.moveTo(obj.x, obj.y);
                this.ctx.lineTo(obj.x + obj.width, obj.y + obj.height);
                this.ctx.stroke();
                break;
                
            case 'rectangle':
                this.ctx.beginPath();
                this.ctx.rect(obj.x, obj.y, obj.width, obj.height);
                this.ctx.stroke();
                break;
                
            case 'circle':
                const radius = Math.sqrt(Math.pow(obj.width, 2) + Math.pow(obj.height, 2)) / 2;
                const centerX = obj.x + obj.width / 2;
                const centerY = obj.y + obj.height / 2;
                this.ctx.beginPath();
                this.ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI);
                this.ctx.stroke();
                break;
                
            case 'text':
                this.ctx.font = `${obj.thickness * 10}px Arial`;
                this.ctx.fillText(obj.text, obj.x, obj.y);
                break;
        }
    }

    redraw() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        // Draw all objects
        this.objects.forEach(obj => {
            this.drawObject(obj);
        });
        
        // Draw selection handles if an object is selected
        if (this.selectedObject) {
            this.drawSelectionHandles();
        }
    }

    handleSelectMouseDown() {
        const clickedObject = this.getObjectAtPosition(this.startX, this.startY);
        
        // Check if clicking on a resize handle
        if (this.selectedObject) {
            const handleIndex = this.getHandleAtPosition(this.startX, this.startY);
            if (handleIndex !== -1) {
                this.isResizing = true;
                this.resizeHandleIndex = handleIndex;
                return;
            }
        }
        
        if (clickedObject) {
            this.selectObject(clickedObject);
        } else {
            this.deselectObject();
        }
    }

    handleSelectMouseMove(currentX, currentY) {
        if (this.selectedObject && !this.isResizing) {
            // Update hover state for resize handles
            const handleIndex = this.getHandleAtPosition(currentX, currentY);
            const handles = document.querySelectorAll('.resize-handle');
            handles.forEach((handle, index) => {
                handle.style.opacity = index === handleIndex ? '1' : '0.5';
            });
        }
    }

    getObjectAtPosition(x, y) {
        // Check objects in reverse order (top to bottom)
        for (let i = this.objects.length - 1; i >= 0; i--) {
            const obj = this.objects[i];
            
            // Expand hit area slightly
            const padding = 10;
            let hitX, hitY, hitWidth, hitHeight;
            
            switch (obj.type) {
                case 'draw':
                case 'erase':
                    // For free draw, check if point is near any segment
                    if (obj.points.length < 2) return null;
                    for (let j = 0; j < obj.points.length - 1; j++) {
                        if (this.isPointNearLine(x, y, obj.points[j], obj.points[j + 1], obj.thickness + padding)) {
                            return obj;
                        }
                    }
                    return null;
                    
                case 'line':
                    hitX = Math.min(obj.x, obj.x + obj.width) - padding;
                    hitY = Math.min(obj.y, obj.y + obj.height) - padding;
                    hitWidth = Math.abs(obj.width) + padding * 2;
                    hitHeight = Math.abs(obj.height) + padding * 2;
                    break;
                    
                case 'rectangle':
                    hitX = obj.x - padding;
                    hitY = obj.y - padding;
                    hitWidth = obj.width + padding * 2;
                    hitHeight = obj.height + padding * 2;
                    break;
                    
                case 'circle':
                    const radius = Math.sqrt(Math.pow(obj.width, 2) + Math.pow(obj.height, 2)) / 2;
                    const centerX = obj.x + obj.width / 2;
                    const centerY = obj.y + obj.height / 2;
                    const distance = Math.sqrt(Math.pow(x - centerX, 2) + Math.pow(y - centerY, 2));
                    return distance <= radius + padding ? obj : null;
                    
                case 'text':
                    hitX = obj.x - padding;
                    hitY = obj.y - obj.thickness * 10 - padding;
                    hitWidth = this.ctx.measureText(obj.text).width + padding * 2;
                    hitHeight = obj.thickness * 10 + padding * 2;
                    break;
            }
            
            if (obj.type !== 'circle' && obj.type !== 'draw' && obj.type !== 'erase') {
                if (x >= hitX && x <= hitX + hitWidth && y >= hitY && y <= hitY + hitHeight) {
                    return obj;
                }
            }
        }
        return null;
    }

    isPointNearLine(px, py, lineStart, lineEnd, tolerance) {
        const A = { x: lineStart.x, y: lineStart.y };
        const B = { x: lineEnd.x, y: lineEnd.y };
        const P = { x: px, y: py };
        
        const l2 = Math.pow(B.x - A.x, 2) + Math.pow(B.y - A.y, 2);
        if (l2 === 0) return Math.sqrt(Math.pow(P.x - A.x, 2) + Math.pow(P.y - A.y, 2)) <= tolerance;
        
        let t = ((P.x - A.x) * (B.x - A.x) + (P.y - A.y) * (B.y - A.y)) / l2;
        t = Math.max(0, Math.min(1, t));
        
        const projection = {
            x: A.x + t * (B.x - A.x),
            y: A.y + t * (B.y - A.y)
        };
        
        const distance = Math.sqrt(Math.pow(P.x - projection.x, 2) + Math.pow(P.y - projection.y, 2));
        return distance <= tolerance;
    }

    getHandleAtPosition(x, y) {
        if (!this.selectedObject) return -1;
        
        const handles = this.getResizeHandles();
        for (let i = 0; i < handles.length; i++) {
            const handle = handles[i];
            const dx = x - handle.x;
            const dy = y - handle.y;
            if (Math.sqrt(dx * dx + dy * dy) <= 10) {
                return i;
            }
        }
        return -1;
    }

    getResizeHandles() {
        if (!this.selectedObject) return [];
        
        const { x, y, width, height } = this.selectedObject;
        return [
            { x, y }, // top-left
            { x: x + width, y }, // top-right
            { x: x + width, y: y + height }, // bottom-right
            { x, y: y + height } // bottom-left
        ];
    }

    drawSelectionHandles() {
        const handles = this.getResizeHandles();
        this.ctx.fillStyle = '#3498db';
        
        handles.forEach(handle => {
            this.ctx.beginPath();
            this.ctx.arc(handle.x, handle.y, 5, 0, 2 * Math.PI);
            this.ctx.fill();
        });
    }

    selectObject(obj) {
        this.deselectObject();
        this.selectedObject = obj;
        this.redraw();
    }

    deselectObject() {
        if (this.selectedObject) {
            this.selectedObject = null;
            this.redraw();
        }
    }

    resizeObject(currentX, currentY) {
        if (!this.selectedObject) return;
        
        const { x, y, width, height } = this.selectedObject;
        let newX = x, newY = y, newWidth = width, newHeight = height;
        
        switch (this.resizeHandleIndex) {
            case 0: // top-left
                newX = currentX;
                newY = currentY;
                newWidth = x + width - currentX;
                newHeight = y + height - currentY;
                break;
            case 1: // top-right
                newY = currentY;
                newX = x;
                newWidth = currentX - x;
                newHeight = y + height - currentY;
                break;
            case 2: // bottom-right
                newX = x;
                newY = y;
                newWidth = currentX - x;
                newHeight = currentY - y;
                break;
            case 3: // bottom-left
                newX = currentX;
                newY = y;
                newWidth = x + width - currentX;
                newHeight = currentY - y;
                break;
        }
        
        // Ensure minimum size
        if (Math.abs(newWidth) > 5 && Math.abs(newHeight) > 5) {
            this.selectedObject.x = newX;
            this.selectedObject.y = newY;
            this.selectedObject.width = newWidth;
            this.selectedObject.height = newHeight;
            this.redraw();
        }
    }

    addTextObject(x, y) {
        console.log('Adding text object at position:', { x, y });
        
        const text = prompt('Enter text:', 'Text');
        if (text) {
            const textObj = {
                type: 'text',
                x: x,
                y: y,
                width: 0,
                height: 0,
                color: this.currentColor,
                thickness: this.currentThickness,
                text: text
            };
            
            console.log('Text object created:', textObj);
            this.objects.push(textObj);
            this.redraw();
            
            // Send text addition event through WebSocket
            this.sendDrawingEvent({
                type: 'textAdd',
                tool: 'text',
                data: {
                    color: textObj.color,
                    thickness: textObj.thickness,
                    x: textObj.x,
                    y: textObj.y,
                    text: textObj.text
                }
            });
        } else {
            console.log('Text input cancelled');
        }
    }

    deleteSelectedObject() {
        if (this.selectedObject) {
            const index = this.objects.indexOf(this.selectedObject);
            if (index > -1) {
                console.log('Deleting object:', {
                    type: this.selectedObject.type,
                    index: index
                });
                
                // Send object deletion event through WebSocket
                this.sendDrawingEvent({
                    type: 'objectDelete',
                    tool: 'select',
                    data: {
                        index: index,
                        objectType: this.selectedObject.type
                    }
                });
                
                this.objects.splice(index, 1);
                this.deselectObject();
                this.redraw();
                
                console.log('Object deleted successfully');
            }
        } else {
            console.log('No object selected for deletion');
        }
    }

    clearBoard() {
        console.log('Clear board requested. Current object count:', this.objects.length);
        
        if (confirm('Are you sure you want to clear the entire board?')) {
            console.log('Board clear confirmed by user');
            
            // Send board clear event through WebSocket
            this.sendDrawingEvent({
                type: 'boardClear',
                tool: 'clear',
                data: {}
            });
            
            this.objects = [];
            this.deselectObject();
            this.redraw();
            console.log('Board cleared successfully');
        } else {
            console.log('Board clear cancelled by user');
        }
    }
}

// Initialize the whiteboard when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new Whiteboard();
});

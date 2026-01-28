const Draw = {
    ctx: null,

    init(canvasId) {
        const can = document.getElementById(canvasId);
        this.ctx = can.getContext('2d');
        // Disable smoothing for pixel art look
        this.ctx.imageSmoothingEnabled = false;
    },

    clear() {
        this.ctx.fillStyle = CONFIG.colors.bg;
        this.ctx.fillRect(0, 0, CONFIG.canvas.width, CONFIG.canvas.height);
    },

    roundRect(x, y, w, h, radius, fill, stroke) {
        this.ctx.beginPath();
        this.ctx.roundRect(x, y, w, h, radius);
        if (fill) {
            this.ctx.fillStyle = fill;
            this.ctx.fill();
        }
        if (stroke) {
            this.ctx.strokeStyle = stroke;
            this.ctx.lineWidth = 4;
            this.ctx.stroke();
        }
    },

    text(str, x, y, size, color, align = 'center') {
        this.ctx.font = `${size}px "Press Start 2P"`;
        this.ctx.fillStyle = color;
        this.ctx.textAlign = align;
        this.ctx.fillText(str, x, y);
    },

    circle(x, y, r, fill, stroke) {
        this.ctx.beginPath();
        this.ctx.arc(x, y, r, 0, Math.PI * 2);
        if (fill) {
            this.ctx.fillStyle = fill;
            this.ctx.fill();
        }
        if (stroke) {
            this.ctx.strokeStyle = stroke;
            this.ctx.lineWidth = 2;
            this.ctx.stroke();
        }
    },

    pixelSprite(sprite, x, y, size, color) {
        if (!sprite) return;
        this.ctx.fillStyle = color;
        for (let r = 0; r < sprite.length; r++) {
            for (let c = 0; c < sprite[r].length; c++) {
                if (sprite[r][c] === 1) {
                    this.ctx.fillRect(x + c * size, y + r * size, size, size);
                }
            }
        }
    }
};

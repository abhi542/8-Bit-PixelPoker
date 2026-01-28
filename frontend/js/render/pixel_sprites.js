const PixelSprites = {
    // 7x7 Bitmaps for Suits
    suits: {
        // Heart (2)
        2: [
            [0, 1, 1, 0, 1, 1, 0],
            [1, 1, 1, 1, 1, 1, 1],
            [1, 1, 1, 1, 1, 1, 1],
            [1, 1, 1, 1, 1, 1, 1],
            [0, 1, 1, 1, 1, 1, 0],
            [0, 0, 1, 1, 1, 0, 0],
            [0, 0, 0, 1, 0, 0, 0]
        ],
        // Diamond (1)
        1: [
            [0, 0, 0, 1, 0, 0, 0],
            [0, 0, 1, 1, 1, 0, 0],
            [0, 1, 1, 1, 1, 1, 0],
            [1, 1, 1, 1, 1, 1, 1],
            [0, 1, 1, 1, 1, 1, 0],
            [0, 0, 1, 1, 1, 0, 0],
            [0, 0, 0, 1, 0, 0, 0]
        ],
        // Club (0)
        0: [
            [0, 0, 0, 0, 0, 0, 0], // Top row empty
            [0, 0, 1, 1, 1, 0, 0], // Top circle
            [0, 0, 1, 1, 1, 0, 0],
            [0, 1, 1, 0, 1, 1, 0], // Side circles top
            [1, 1, 1, 1, 1, 1, 1], // Middle
            [0, 0, 0, 1, 0, 0, 0], // Stem
            [0, 0, 1, 1, 1, 0, 0]  // Base
        ],
        // Spade (3)
        3: [
            [0, 0, 0, 1, 0, 0, 0], // Tip
            [0, 0, 1, 1, 1, 0, 0],
            [0, 1, 1, 1, 1, 1, 0],
            [1, 1, 1, 1, 1, 1, 1],
            [1, 1, 1, 1, 1, 1, 1],
            [1, 0, 0, 1, 0, 0, 1], // Legs
            [0, 0, 0, 1, 0, 0, 0]  // Stem base
        ]
    },
    // We can also define Rank patterns if custom font is overkill, 
    // but relying on a Pixel Font file is cleaner for text.
    // For Card Back: 4x4 tiling pattern
    cardBackPattern: [
        [1, 0, 1, 0],
        [0, 1, 0, 1],
        [1, 0, 1, 0],
        [0, 1, 0, 1]
    ]
};

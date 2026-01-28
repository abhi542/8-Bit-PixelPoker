const CardRenderer = {
    drawCard(card, x, y, scale = 1) {
        const w = 40 * scale;
        const h = 56 * scale;

        // Shadow
        Draw.roundRect(x + 2, y + 2, w, h, 4, '#000');

        // Base
        Draw.roundRect(x, y, w, h, 4, CONFIG.colors.cardWhite);

        if (!card) {
            // Back of card
            Draw.roundRect(x + 4, y + 4, w - 8, h - 8, 2, CONFIG.colors.cardRed);
            return;
        }

        // Parse Code (e.g. "As", "Td")
        // Check local struct: Backend "Card" is {Suit:int, Rank:int}.
        // But JSON marshals to {suit:X, rank:Y}? Or did we add custom marshal?
        // Checking internal/poker/card.go: `json:"suit"` `json:"rank"`.
        // Suits: 0=Clubs, 1=Diamonds, 2=Hearts, 3=Spades
        // Ranks: 2=2... 14=Ace

        const suits = ['♣', '♦', '♥', '♠'];
        const colors = [CONFIG.colors.cardBlack, CONFIG.colors.cardRed, CONFIG.colors.cardRed, CONFIG.colors.cardBlack];
        const ranks = { 11: 'J', 12: 'Q', 13: 'K', 14: 'A', 10: '10' };

        const sIdx = card.suit; // 0-3
        const rVal = card.rank;

        const rankStr = ranks[rVal] || rVal;
        const suitChar = suits[sIdx];
        const color = colors[sIdx];

        Draw.text(rankStr, x + w / 2, y + h / 2 - 5, 12 * scale, color, 'center');
        Draw.text(suitChar, x + w / 2, y + h / 2 + 10, 14 * scale, color, 'center');
    }
};

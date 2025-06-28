export default class ColorService {
    constructor() {
        // Base color palette (extend as needed)
        this.BASE_COLORS = [
            '#FF5252', '#4CAF50', '#2196F3', '#FFC107', // Vibrant primary
            '#9C27B0', '#00BCD4', '#FF9800', '#795548', // Secondary options
            '#E91E63', '#8BC34A', '#3F51B5', '#FF5722'  // More alternatives
        ];

        // Track assigned colors
        this.assignedColors = new Map();
    }

    // Get or create color for participant
    getParticipantColor(participantID) {
        // Handle empty participant ID
        if (!participantID) {
            return '#ffffff';
        }

        // Return existing if already assigned
        if (this.assignedColors.has(participantID)) {
            return this.assignedColors.get(participantID);
        }

        // Assign new color
        const color = this.generateColor();
        this.assignedColors.set(participantID, color);
        return color;
    }

    // Generate color based on current assignment count
    generateColor() {
        // Cycle through base palette
        return this.BASE_COLORS[this.assignedColors.size % this.BASE_COLORS.length];
    }

    // Reset all color assignments
    reset() {
        this.assignedColors.clear();
    }
}
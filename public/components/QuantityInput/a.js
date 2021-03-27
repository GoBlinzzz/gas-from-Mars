export default class QuantityInput {
    constructor(self, decreaseText, increaseText) {
        this.input = document.createElement('input');
        this.input.value = '1';
        this.input.type = 'number';
        this.input.name = 'quantity';
        this.input.pattern = '[0-9]+';
        this.decreaseText = decreaseText || 'Decrease quantity';
        this.increaseText = increaseText || 'Increase quantity';

        function Button(text, className) {
            this.button = document.createElement('button');
            this.button.type = 'button';
            this.button.innerHTML = text;
            this.button.title = text;
            this.button.classList.add(className);
            return this.button;
        }

        this.subtract = new Button(this.decreaseText, 'sub');
        this.add = new Button(this.increaseText, 'add');
        this.subtract.addEventListener('click', () => this.change_quantity(-1));
        this.add.addEventListener('click', () => this.change_quantity(1));
        self.appendChild(this.subtract);
        self.appendChild(this.input);
        self.appendChild(this.add);
    }

    change_quantity(change) {
        let quantity = Number(this.input.value);
        if (isNaN(quantity)) quantity = 1;
        quantity += change;
        quantity = Math.max(quantity, 1);
        this.input.value = String(quantity);
    }
}

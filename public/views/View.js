 export class View {
    constructor(parent, templatePath) {
        this._parent = parent;
        this._name = templatePath;
    }

     render(data = {}) {
         this._parent.innerHTML = Handlebars.templates[this._name](data);
     }
 }
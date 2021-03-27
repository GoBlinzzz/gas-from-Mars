import {View} from '../View.js';
import {MainPageView} from "../MainPageView/MainPageView.js";
import http from "../../modules/http.js";

export class ConfigPageView extends View {
    constructor(parent) {
        super(parent, 'views/ConfigPageView/ConfigPageView');
        this._parent = parent;
    }

    render(data = {}) {
        super.render(data);
        const submitButton = this._parent.getElementsByClassName('submit-config__button').item(0);
        submitButton.addEventListener('click', () => {
            const input = this._parent.getElementsByClassName('file-input__input')[0];

            const file = input.files[0];

            let json = {};

            const reader = new FileReader();

            reader.addEventListener('load', (event) => {
                const result = event.target.result;
                json = JSON.parse(result.toString());

                http.post('/config', json)
                    .then((res) => {
                        if (res.status === 200) {
                            const app = document.getElementById('app');
                            const mainPage = new MainPageView(app);
                            http.get('/simulation/1')
                                .then((res) => {
                                    if (res.status === 200) {
                                        mainPage.render(res.body);
                                    }
                                });
                        }
                    });
            });

            reader.readAsBinaryString(file);
        });
    }
}

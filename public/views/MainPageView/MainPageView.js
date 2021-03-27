import {View} from "../View.js";
import http from "../../modules/http.js";

export class MainPageView extends View {
    constructor(parent) {
        super(parent, 'views/MainPageView/MainPageView');
    }

    render(data = {}) {
        super.render(data);

        this._interval = data.interval * 1000;
        this.processYear();
    }

    reRender(data) {

    }

    sendMonthRequest(month) {
        http.get(`/simulation/${month}`).then((res) => {
            if (res.status === 200) {
                this.reRender(res.body);
            }
        });
    }

    processYear() {
        for (let i = 2; i <= 12; i++) {
            setTimeout(() => {
                this.sendMonthRequest(i)
            }, this._interval);
        }
    }
}
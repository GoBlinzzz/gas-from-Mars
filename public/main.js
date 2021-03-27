import {MainPageView} from './views/MainPageView/MainPageView.js';
import {ConfigPageView} from "./views/ConfigPageView/ConfigPageView.js";

const app = document.getElementById('app');

const mainPage = new MainPageView(app);
mainPage.render();

const configPage = new ConfigPageView(app);
configPage.render();

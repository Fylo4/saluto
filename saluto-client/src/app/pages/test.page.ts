import { Component, inject, signal } from "@angular/core";
import { ApiService } from "../services/api.service";

@Component({
    selector: 'app-test-page',
    templateUrl: './test.page.html',
    styleUrl: './test.page.scss',
    standalone: true,
    imports: []
})
export class TestPage {
    private api = inject(ApiService);

    serverResponse = signal("...");

    testServer() {
        this.api.testServer().subscribe(v => {
            this.serverResponse.set(v.message);
        })
    }
}
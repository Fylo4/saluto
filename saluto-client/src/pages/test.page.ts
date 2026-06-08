import { Component, inject, OnInit, signal } from "@angular/core";
import { ApiService } from "../services/api.service";
import { FormsModule } from "@angular/forms";
import { API_Post_Get_Return } from "../services/api.service.types";
import { DatePipe } from "@angular/common";

@Component({
    selector: 'app-test-page',
    templateUrl: './test.page.html',
    styleUrl: './test.page.scss',
    standalone: true,
    imports: [FormsModule, DatePipe]
})
export class TestPage implements OnInit {
    private api = inject(ApiService);

    testServerResponse = signal("");
    messageToPost = signal("");
    displayName = signal("");
    errMessage = signal("");
    allMessages = signal<API_Post_Get_Return[]>([]);

    ngOnInit() {
        this.btnGetMessages();
    }
    
    testServer() {
        this.api.testServer().subscribe(v => {
            this.testServerResponse.set(v.message);
            setTimeout(() => this.testServerResponse.set(''), 2000)
        })
    }

    btnPostMessage() {
        if (!this.validateMessage()) return;

        const messageInput = this.messageToPost();
        const displayNameInput = this.displayName();
        this.api.postMessage({
            message: messageInput,
            displayName: displayNameInput
        }).subscribe(v => {
            if (v !== true) return;

            this.messageToPost.set("");
            this.errMessage.set("");
            const newMessage: API_Post_Get_Return = {
                ID: -1,
                Body: messageInput,
                Displayname: displayNameInput,
                Timeposted: new Date()
            }
            this.allMessages.set(
                [ newMessage, ...this.allMessages()]
            )
        })
    }
    btnGetMessages() {
        this.api.getMessages().subscribe(v => {
            this.allMessages.set(v);
        });
    }

    private validateMessage() {
        if (this.displayName().trim() == '') {
            this.errMessage.set("Display name cannot be blank");
            return false;
        }
        if (this.messageToPost().trim() == '') {
            this.errMessage.set("Message cannot be blank");
            return false;
        }
        return true;
    }
}
import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { API_Post_Create_Input, API_Post_Get_Return } from "./api.service.types";
import { map } from "rxjs";
import { ConfigService } from "./config.service";

@Injectable({
    providedIn: 'root'  
})
export class ApiService {
    private http = inject(HttpClient);
    private cfg = inject(ConfigService);

    public testServer() {
        return this.http.get<{message: string}>(this.cfg.APIRoot+"test")
    }

    postMessage(input: API_Post_Create_Input) {
        return this.http.post<boolean>(this.cfg.APIRoot+"post", input)
    }

    getMessages() {
        return this.http.get<API_Post_Get_Return[]>(this.cfg.APIRoot+"posts")
            .pipe(map(messages => {
                return messages.map(message => ({
                    ...message,
                    Timeposted: new Date(message.Timeposted)
                }))
            }));
    }
}
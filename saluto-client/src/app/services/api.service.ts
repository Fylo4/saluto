import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";

@Injectable({
    providedIn: 'root'  
})
export class ApiService {
    private http = inject(HttpClient);

    public testServer() {
        return this.http.get<{message: string}>("/api/")
    }
}
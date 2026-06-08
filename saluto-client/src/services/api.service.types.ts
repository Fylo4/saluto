export type API_Post_Create_Input = {
    displayName: string;
    message: string;
}
export type API_Post_Get_Return = {
    ID: number;
    Body: string;
    Displayname: string;
    Timeposted: Date;
}
import { createContext } from "react";

export const FilesContext = createContext({
    files: [],
    addFiles: () => {},
    deleteFile: () => {},
});

export const AuthContext = createContext({
    changeAuthStatus: () => {},
})

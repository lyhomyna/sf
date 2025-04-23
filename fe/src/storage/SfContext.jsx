import { createContext } from "react";

export const FilesContext = createContext({
    filenames: [],
    addFilenames: () => {},
    deleteFilename: () => {},
});

export const AuthContext = createContext({
    changeAuthStatus: () => {},
})

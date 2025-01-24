import { createContext } from "react";

export const FilesContext = createContext({
    files: [],
    addFilenames: () => {},
    addFilename: () => {}
});

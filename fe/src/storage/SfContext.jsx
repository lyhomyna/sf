import { createContext } from "react";

export const DirItemsContext = createContext({
    dirItems: [],
    addDirItems: () => {},
    deleteDirItem: () => {},
});

export const AuthContext = createContext({
    changeAuthStatus: () => {},
})

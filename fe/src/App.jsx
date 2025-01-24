import { useState } from "react";

import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";
import { FilesContext } from "./storage/FilesContext.jsx";

export default function App() {
    const [filenames, setFilenames] = useState([]);

    const email="supercoolemail@super.mail" 
    const imageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"
    
    const addFilenames = (filenames) => {
	if (Array.isArray(filenames)) {
	    setFilenames((oldFilenames) => [...oldFilenames, ...filenames] );
	}
    };

    const addFilename = (filename) => {
	setFilenames(oldFilenames => [...oldFilenames, filename])
    };

    return (
    <FilesContext value={ { filenames: filenames, addFilenames: addFilenames, addFilename: addFilename } }>
	<div className="p-2" >
	    <TopBar email={email} imageURL={imageURL} />
	    <FileList />
	</div>
    </FilesContext>
    );
}

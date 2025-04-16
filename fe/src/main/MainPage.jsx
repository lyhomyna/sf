import { useState } from "react";

import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";
import DragAndDrop from "./components/DragAndDrop.jsx";
import { FilesContext } from "storage/SfContext.jsx";

export default function MainPage() {
    const [filenames, setFilenames] = useState([]);

    const addFilenames = ({ filenames, rewrite=false }) => {
	if (rewrite) {
	    setFilenames(filenames);
	} else {
	    setFilenames(oldFilenames => [...oldFilenames, ...filenames]);
	}
    };

    const deleteFilename = (filename) => {
	setFilenames(oldFilenames => { 
	    let withoutFilename;
	    oldFilenames.forEach((_, i) => {
		if (oldFilenames[i] === filename) {
		    withoutFilename = [...oldFilenames.slice(0, i), ...oldFilenames.slice(i+1, oldFilenames.length)]
		    return
		}
	    });

	    return withoutFilename;
	});
    };

    return <FilesContext value={ { 
	filenames: filenames, 
	addFilenames: addFilenames, 
	deleteFilename: deleteFilename, 
    }}>
	<div className="p-2" >
	    <TopBar />
	    <FileList />
	    <DragAndDrop /> 
	</div>
    </FilesContext>
}

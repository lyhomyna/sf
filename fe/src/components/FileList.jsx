import { useState } from 'react';

import FileItem from "./FileItem.jsx";

export default function FileList() {
    const [filenames, setFilenames] = useState([]); 

    fetch("http://localhost:8080/filenames")

    return <ul className="flex flex-col justify-start w-max">
	<FileItem fullFilename="supercoolfile1.txt"/>
	<FileItem fullFilename="supercoolfilee2.txt"/>
	<FileItem fullFilename="supercoolfileee3.txt"/>
    </ul>
}

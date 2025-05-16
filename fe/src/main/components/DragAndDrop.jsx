import { useEffect, useState, useCallback, useRef } from 'react';
import { useDropzone } from 'react-dropzone';
import { useUploadFile } from "hooks/useUploadFile.js"; 

function DragAndDropSection() {
    const { getRootProps, getInputProps } = useDropzone();

    return (
    <div
      {...getRootProps()}
      className="drag-n-drop transition-all fixed inset-0 z-50 flex min-h-screen items-center justify-center bg-gray-900/70 backdrop-blur"
    >
      <div
	aria-hidden="true"
	className="absolute inset-y-16 inset-x-0 w-16 rounded-full rotate-45 bg-gradient-to-b from-pink-500 to-purple-600 blur-3xl mx-auto scale-y-150 opacity-75"
      />
      <input className="w-full h-full" {...getInputProps()} />
      <h1 className="relative z-10 text-4xl font-bold text-white">
	Drop file or files here
      </h1>
    </div>
    );
}

export default function DragAndDrop() {
    const { uploadFile } = useUploadFile();
    const [isDragging, setIsDragging] = useState(false);
    let dragCounter = useRef(0);

    const onDrop = useCallback(async (acceptedFiles) => {
	setIsDragging(false);
	await Promise.all(acceptedFiles.map(file => uploadFile(file)));
    }, [uploadFile]);

    useEffect(() => {
    const handleDragEnter = (e) => {
	e.preventDefault();

	// show dropzone only when files are being dragged 
	if (!e.dataTransfer.types.includes('Files')) {
	    return;
	}

	dragCounter.current++;
	setIsDragging(true);
    };

    const handleDragLeave = (e) => {
	e.preventDefault();

	// show dropzone only when files are being dragged 
	if (!e.dataTransfer.types.includes('Files')) {
	    alert("Isn't a file")
	    return; 
	}

	dragCounter.current--;
	if (dragCounter === 0) {
	setIsDragging(false);
	}
    };

    const handleDrop = (e) => {
	e.preventDefault();

	// show dropzone only when files are being dragged 
	if (!e.dataTransfer.types.includes('Files')) {
	    alert("Isn't a file")
	    return; 
	}

	dragCounter.current = 0;
	setIsDragging(false);
    };

    const preventDefaults = (e) => {
	e.preventDefault();
	e.stopPropagation();
    };

    window.addEventListener('drop', handleDrop);
    window.addEventListener('dragenter', handleDragEnter);
    window.addEventListener('dragleave', handleDragLeave);
    window.addEventListener('dragover', preventDefaults);

    return () => {
      window.removeEventListener('drop', handleDrop);
      window.removeEventListener('dragenter', handleDragEnter);
      window.removeEventListener('dragleave', handleDragLeave);
      window.removeEventListener('dragover', preventDefaults);
    };
    }, []);

    return isDragging ? <DragAndDropSection onDrop={onDrop} /> : null;
}

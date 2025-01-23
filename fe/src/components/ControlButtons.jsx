import Button from "./Button.jsx";

export default function ControlButtons({...props}) {
    return (<div {...props}>
	<Button text="Log in" />
	<Button text="Cnange password" />
	<Button text="Upload file" onClick={uploadFile} />
    </div>);
}

function uploadFile() {
    const input = document.createElement('input');
    input.type = 'file';

    input.onchange = async e => {
	const file = e.target.files[0];
	const formData = new FormData();
	formData.append("file", file);
	
	try {
	    const response = await fetch("http://localhost:8080/save", {
		method: "POST", 
		body: formData,
	    });

	    if (response.ok) {
		alert("File uploaded successfuly!");
	    } else {
		alert("Failed to upload file.");
	    }
	} catch (err) {
	    console.error("Error uploading file:", err)
	    alert("An error occured. Try again.")
	}
    }

    input.click();
}

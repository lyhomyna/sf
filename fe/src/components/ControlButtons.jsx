import Button from "./Button.jsx";

export default function ControlButtons({...props}) {
    return (<div {...props}>
	<Button text="Log in" />
	<Button text="Cnange password" />
	<Button text="Upload file" />
    </div>);
}

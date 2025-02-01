export default function Button({ text, className, ...props }) {
    return (<button className={"rounded-lg text-neutral-400 py-2 px-3 hover:text-neutral-300 duration-300 "+className} {...props}>
	{text}
    </button>);
}


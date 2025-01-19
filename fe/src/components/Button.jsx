export default function Button({ text, className, ...props }) {
    return (<button className={className+" rounded-lg text-neutral-400 bg-neutral-700 py-2 px-3 hover:text-neutral-300 duration-300"} {...props}>
	{text}
    </button>
    );
}


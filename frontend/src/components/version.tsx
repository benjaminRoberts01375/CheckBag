const Version = () => {
	return (
		<>
			<div className="version">
				<p>CheckBag: </p>
				{__CHECKBAG_VERSION__ == "dev" ? (
					<p>Development Build</p>
				) : (
					<a
						href={`https://github.com/benjaminRoberts01375/CheckBag/releases/tag/${__CHECKBAG_VERSION__}`}
					>
						{__CHECKBAG_VERSION__}
					</a>
				)}
			</div>
		</>
	);
};

export default Version;

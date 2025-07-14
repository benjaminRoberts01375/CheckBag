import "../styles.css";
import Service from "../types/service.tsx";
import { useState } from "react";
import { useList } from "../context-hook";

interface ServiceAddScreenProps {
	service: Service | undefined;
}

const ServiceEdit = ({ service }: ServiceAddScreenProps) => {
	const [name, setName] = useState<string>(service ? service.title : "");
	const [internalAddress, setInternalAddress] = useState<string>(
		service ? service.internal_address : "",
	);
	const [externalAddress, setExternalAddress] = useState<string[]>(
		service ? service.external_address : [""],
	);
	const { serviceAdd } = useList();
	function createService() {
		console.log("Creating service");
		serviceAdd(new Service(internalAddress, externalAddress, name));
		setName("");
		setInternalAddress("");
		setExternalAddress([""]);
	}

	return (
		<tr>
			<td>
				<input
					type="text"
					placeholder="Name"
					value={name}
					onChange={e => setName(e.target.value)}
				/>
			</td>
			<td>
				<input
					type="text"
					placeholder="External Address"
					value={externalAddress}
					onChange={e => setExternalAddress([e.target.value])}
				/>
			</td>
			<td>
				<input
					type="text"
					placeholder="Internal Address"
					value={internalAddress}
					onChange={e => setInternalAddress(e.target.value)}
				/>
			</td>
			<td>
				{service ? (
					<button
						onClick={() => {
							console.log("Service exists");
						}}
					>
						Exists
					</button>
				) : (
					<button
						onClick={() => createService()}
						disabled={name === "" || internalAddress === "" || externalAddress[0] === ""}
					>
						Add
					</button>
				)}
			</td>
		</tr>
	);
};

export default ServiceEdit;

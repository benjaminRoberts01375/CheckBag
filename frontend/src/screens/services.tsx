import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useState, useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Services</title>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				{services.map(service => (
					<ServiceEntry service={service} key={service.clientID} />
				))}
				<button className="submit" id={ServicesStyles["add-service-button"]}>
					Add Service
				</button>
			</div>
		</div>
	);
};

export default ServicesScreen;

interface ServiceListEntryProps {
	service: Service;
}

const ServiceEntry = ({ service }: ServiceListEntryProps) => {
	const dialogRef = useRef<HTMLDialogElement | null>(null);

	function openDialog(): void {
		console.log("Opening dialog");
		dialogRef.current?.showModal();
	}

	function closeDialog(): void {
		dialogRef.current?.close();
	}

	return (
		<div id={ServicesStyles["service-container"]}>
			<dialog ref={dialogRef}>
				<EditService service={service} />
			</dialog>
			<h2>{service ? service.title : "Untitled Service"}</h2>
			<div className={ServicesStyles["connection-info"]}>
				<div id={ServicesStyles["service-endpoints"]}>
					{service?.external_address.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service?.internal_address ? <ServiceStatus address={service.internal_address} /> : null}
				<button onClick={() => openDialog()}>
					<CgMoreVerticalAlt className="icon" id={ServicesStyles["menu-icon"]} />
				</button>
			</div>
		</div>
	);
};

interface ServiceURLProps {
	address: string;
}

const ServiceStatus = ({ address }: ServiceURLProps) => {
	return <p className={ServicesStyles["service-status"]}>{address}</p>;
};

const EditService = ({ service }: ServiceListEntryProps) => {
	const [title, setTitle] = useState(service.title || "Untitled Service");
	const [incomingAddress, setIncomingAddress] = useState(service.external_address[0] || "");
	const [outgoingAddress, setOutgoingAddress] = useState(service.internal_address || "");

	return (
		<form id={ServicesStyles["edit-service-container"]}>
			<h1>Editing "{title}"</h1>
			<input
				type="text"
				placeholder="Service Name"
				value={title}
				onChange={e => setTitle(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<p>From:</p>
			<input
				type="text"
				placeholder="Forward Address"
				value={incomingAddress}
				onChange={e => setIncomingAddress(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<p>To:</p>
			<input
				type="text"
				placeholder="Forward Address"
				value={outgoingAddress}
				onChange={e => setOutgoingAddress(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<div id={ServicesStyles["submission-buttons"]}>
				<button className="delete">Delete</button>
				<button className="cancel">Cancel</button>
				<button className="submit">Submit</button>
			</div>
		</form>
	);
};

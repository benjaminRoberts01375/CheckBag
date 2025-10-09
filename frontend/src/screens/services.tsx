import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useState, useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";
import { CommunicationProtocols, CommunicationProtocol } from "../types/strings";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);
	const dialogRef = useRef<HTMLDialogElement | null>(null);
	const [isDialogOpen, setIsDialogOpen] = useState(false);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Services</title>
			<dialog ref={dialogRef}>
				{isDialogOpen && (
					<EditService
						service={undefined}
						finish={() => {
							dialogRef.current?.close();
							setIsDialogOpen(false);
						}}
					/>
				)}
			</dialog>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				{services.map(service => (
					<ServiceEntry servicePass={service} key={service.clientID} />
				))}
				<button
					className="submit"
					id={ServicesStyles["add-service-button"]}
					onClick={() => {
						setIsDialogOpen(true);
						dialogRef.current?.showModal();
					}}
				>
					Add Service
				</button>
			</div>
		</div>
	);
};

export default ServicesScreen;

interface ServiceListEntryProps {
	servicePass: Service;
}

const ServiceEntry = ({ servicePass }: ServiceListEntryProps) => {
	const dialogRef = useRef<HTMLDialogElement | null>(null);
	const [service, _] = useState<Service>(servicePass);
	const [isDialogOpen, setIsDialogOpen] = useState(false);

	return (
		<div id={ServicesStyles["service-container"]}>
			<dialog ref={dialogRef}>
				{isDialogOpen && (
					<EditService
						service={service}
						finish={() => {
							dialogRef.current?.close();
							setIsDialogOpen(false);
						}}
					/>
				)}
			</dialog>
			<h2>{service.title}</h2>
			<div className={ServicesStyles["connection-info"]}>
				<div id={ServicesStyles["service-endpoints"]}>
					{service.external_address.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service.internal_address ? <ServiceStatus address={service.internal_address} /> : null}
				<button
					onClick={() => {
						setIsDialogOpen(true);
						dialogRef.current?.showModal();
					}}
				>
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

interface EditServiceProps {
	service: Service | undefined;
	finish: () => void;
}

const EditService = ({ service, finish }: EditServiceProps) => {
	const { serviceAdd, serviceDelete, serviceUpdate } = useList();
	const [title, setTitle] = useState(service?.title ?? "");
	const [incomingAddresses, setIncomingAddress] = useState(service?.external_address ?? []);

	const [outgoingProtocol, setOutgoingProtocol] = useState<CommunicationProtocol>(
		service?.internal_address?.startsWith("https") ? "https" : "http",
	);
	const [outgoingDomain, setOutgoingDomain] = useState(
		service?.internal_address.split(":")[0] ?? "", // ex. TODO: Get entry 1 from split and remove first two characters
	);
	const [outgoingPort, setOutgoingPort] = useState(() => {
		console.log("Checking port:", service?.internal_address.split(":")[1]);
		return service?.internal_address.split(":")[1] ?? "80";
	}); // TODO: Get entry 2 from split

	function submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		service == undefined ? createService() : updateService(service);
		finish();
	}

	function createService() {
		console.log("Adding service:", outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort);
		serviceAdd(
			new Service(
				outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort,
				incomingAddresses,
				title,
			),
		);
	}

	function updateService(service: Service) {
		service.title = title;
		service.external_address = incomingAddresses;
		service.internal_address = outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort;
		serviceUpdate(service);
	}

	function cancel(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Cancelling");
		finish();
	}

	function deleteService(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Deleting");
		if (service) {
			serviceDelete(service.clientID);
		}
		finish();
	}

	return (
		<form id={ServicesStyles["edit-service-container"]}>
			<div id={ServicesStyles["edit-service-header"]}>
				<h1>Editing</h1>
				<input
					type="text"
					placeholder="Service Name"
					value={title}
					onChange={e => setTitle(e.target.value)}
					className={ServicesStyles["input"]}
				/>
			</div>

			<div className={ServicesStyles["url-container"]}>
				<h3>From:</h3>
				<input
					type="url"
					autoComplete="off"
					placeholder="Forward Address"
					value={incomingAddresses}
					onChange={e => setIncomingAddress([e.target.value])}
					className={ServicesStyles["input"]}
				/>
			</div>
			<div className={ServicesStyles["url-container"]}>
				<h3>To:</h3>
				{/* <select
					value={outgoingProtocol}
					onChange={e => setOutgoingProtocol(e.target.value as CommunicationProtocol)}
					className={ServicesStyles["input"]}
				>
					{CommunicationProtocols.map(protocol => (
						<option value={protocol} key={protocol}>
							{protocol}
						</option>
					))}
				</select>
				<p>://</p> */}
				<input
					type="url"
					autoComplete="off"
					placeholder="Forward Domain"
					value={outgoingDomain}
					onChange={e => setOutgoingDomain(e.target.value)}
					className={ServicesStyles["input"]}
				/>
				<p>:</p>
				<input
					type="number"
					autoComplete="off"
					placeholder="Port"
					value={outgoingPort}
					onChange={e => {
						var newPort = Number(e.target.value);
						if ((newPort > 0 && newPort <= 65535) || e.target.value == "") {
							setOutgoingPort(e.target.value);
						}
					}}
					min={1}
					max={65535}
					className={`${ServicesStyles["input"]} ${ServicesStyles["port"]}`}
				/>
			</div>
			<div id={ServicesStyles["submission-buttons"]}>
				{service ? (
					<button
						className="delete"
						role="delete"
						onClick={e => {
							deleteService(e);
						}}
					>
						Delete
					</button>
				) : null}

				<button
					className="cancel"
					role="cancel"
					onClick={e => {
						cancel(e);
					}}
				>
					Cancel
				</button>
				<button
					className="submit"
					role="submit"
					onClick={e => {
						submit(e);
					}}
				>
					{service == undefined ? "Create" : "Save"}
				</button>
			</div>
		</form>
	);
};

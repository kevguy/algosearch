import React from 'react';
import axios from 'axios';
import moment from 'moment';
import './index.css';
import { NavLink } from 'react-router-dom';
import Layout from '../../components/layout';
import Breadcrumbs from '../../components/breadcrumbs';
import Load from '../../components/tableloading';
import AlgoIcon from '../../components/algoicon';
import {formatValue, siteName} from '../../constants';

class Transaction extends React.Component {
	constructor() {
		super();

		this.state = {
			txid: 0,
			transaction: [],
			loading: true,
		}
	}

	getTransaction = txid => {
		axios({
			method: 'get',
			url: `${siteName}/transactionservice/${txid}`
		}).then(response => {
			this.setState({transaction: response.data, loading: false});
		}).catch(error => {
			console.log("Exception when retrieving transaction details: " + error);
		})
	};

	componentDidMount() {
		const { txid } = this.props.match.params;
		this.setState({txid: txid});
		document.title=`AlgoSearch | Transaction ${txid}`;
		this.getTransaction(txid);
	}

	render() {
		return (
			<Layout>
				<Breadcrumbs
					name={`Transaction Details`}
					parentLink="/transactions"
					parentLinkName="Transactions"
					currentLinkName={`Transaction Details`}
				/>
				<div className="block-table">
					<span>Transaction Details</span>
					<div>
						<table cellSpacing="0">
							<thead>
								<tr>
									<th>Identifier</th>
									<th>Value</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td>ID</td>
									<td>{this.state.loading ? <Load /> : this.state.transaction.transaction.id}</td>
								</tr>
								<tr>
									<td>Round</td>
									<td>{this.state.loading ? <Load /> : <NavLink to={`/block/${this.state.transaction.transaction['comfirmed-round']}`}>{this.state.transaction.transaction['confirmed-round']}</NavLink>}</td>
								</tr>
								<tr>
									<td>Type</td>
									<td>{this.state.loading ? <Load /> : <span className="type noselect">{this.state.transaction.transaction['tx-type']}</span>}</td>
								</tr>
								<tr>
									<td>Sender</td>
									<td>{this.state.loading ? <Load /> : <NavLink to={`/address/${this.state.transaction.transaction.sender}`}>{this.state.transaction.transaction.sender}</NavLink>}</td>
								</tr>
								<tr>
									<td>Receiver</td>
									<td>{this.state.loading ? <Load /> : <NavLink to={`/address/${this.state.transaction.transaction['payment-transaction'].receiver}`}>{this.state.transaction.transaction['payment-transaction'].receiver}</NavLink>}</td>
								</tr>
								<tr>
									<td>Amount</td>
									<td>{this.state.loading ? <Load /> : (
										<div className="tx-hasicon">
											{formatValue(this.state.transaction.transaction['payment-transaction'].amount / 1000000)}
											<AlgoIcon />
										</div>
									)}</td>
								</tr>
								<tr>
									<td>Fee</td>
									<td>{this.state.loading ? <Load /> : (
										<div className="tx-hasicon">
											{formatValue(this.state.transaction.transaction.fee / 1000000)}
											<AlgoIcon />
										</div>
									)}</td>
								</tr>
								<tr>
									<td>First round</td>
									<td>{this.state.loading ? <Load /> : <NavLink to={`/block/${this.state.transaction.transaction["first-valid"]}`}>{this.state.transaction.transaction["first-valid"]}</NavLink>}</td>
								</tr>
								<tr>
									<td>Last round</td>
									<td>{this.state.loading ? <Load /> : <NavLink to={`/block/${this.state.transaction.transaction["last-valid"]}`}>{this.state.transaction.transaction["last-valid"]}</NavLink>}</td>
								</tr>
								<tr>
									<td>Timestamp</td>
									<td>{this.state.loading ? <Load /> : moment.unix(this.state.transaction.timestamp).format("LLLL")}</td>
								</tr>
								<tr>
									<td>Note</td>
									<td>
										{this.state.loading ? <Load /> : (
											<div>
												{this.state.transaction.transaction.note && this.state.transaction.transaction.note !== '' ? (
													<div>
														<div>
															<span>Base 64:</span>
															<textarea defaultValue={this.state.transaction.transaction.note} readOnly></textarea>
														</div>
														<div>
															<span>Converted:</span>
															<textarea defaultValue={atob(this.state.transaction.transaction.note)} readOnly></textarea>
														</div>
													</div>
												) : null}
											</div>
										)}
									</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
				<div className="block-table">
					<span>Miscellaneous Details</span>
					<div>
						<table cellSpacing="0">
							<thead>
								<tr>
									<th>Identifier</th>
									<th>Value</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td>From rewards</td>
									<td>{this.state.loading ? <Load /> : (
										<div className="tx-hasicon">
											{formatValue(this.state.transaction.transaction['sender-rewards'] / 1000000)}
											<AlgoIcon />
										</div>
									)}</td>
								</tr>
								<tr>
									<td>To rewards</td>
									<td>{this.state.loading ? <Load /> : (
										<div className="tx-hasicon">
											{formatValue(this.state.transaction.transaction['receiver-rewards'] / 1000000)}
											<AlgoIcon />
										</div>
									)}</td>
								</tr>
								<tr>
									<td>Genesis ID</td>
									<td>{this.state.loading ? <Load /> : this.state.transaction.transaction['genesis-id']}</td>
								</tr>
								<tr>
									<td>Genesis hash</td>
									<td>{this.state.loading ? <Load /> : this.state.transaction.transaction['genesis-hash']}</td>
								</tr>
							</tbody>
						</table>
					</div>
				</div>
			</Layout>
		);
	}
}

export default Transaction;

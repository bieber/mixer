/*
 * Copyright 2015, Robert Bieber
 *
 * This file is part of mixer.
 *
 * mixer is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mixer is distributed in the hope that it will be useful,
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mixer.  If not, see <http://www.gnu.org/licenses/>.
 */

import React from 'react/addons';

import ListItem from '../components/ListItem.js';
import AddableItemFooter from '../components/AddableItemFooter.js';
import RemovableItemFooter from '../components/RemovableItemFooter.js';

var {update} = React.addons;

export default class Composer extends React.Component {
	constructor(props, context) {
		super(props, context);
		this.state = {
			sourceLists: [],
			destList: null,
		};
	}

	onAddToSource(list, event) {
		event.preventDefault();
		this.setState({
			sourceLists: update(this.state.sourceLists, {$push: [list]}),
		});
	}

	onRemoveFromSource(index, event) {
		event.preventDefault();
		this.setState({
			sourceLists: update(
				this.state.sourceLists,
				{$splice: [[index, 1]]}
			),
		});
	}

	onSetAsDest(list, event) {
		event.preventDefault();
		this.setState({destList: list});
	}

	onRemoveDest(event) {
		event.preventDefault();
		this.setState({destList: null});
	}

	render() {
		if (!this.props.playlists) {
			return <h1>Loading playlists...</h1>;
		}

		var usedIDs = {};
		for (var i = 0; i < this.state.sourceLists.length; i++) {
			usedIDs[this.state.sourceLists[i].id] = true;
		}
		if (this.state.destList) {
			usedIDs[this.state.destList.id] = true;
		}

		var freeLists = this.props.playlists
			.filter(
				(l, i, ls) => !(l.id in usedIDs)
			).map(
				(l, i, ls) => (
					<ListItem key={l.id} list={l}>
						<AddableItemFooter
							onAddToSource={this.onAddToSource.bind(this, l)}
							onSetAsDest={this.onSetAsDest.bind(this, l)}
							isOwned={this.props.userID === l.owner.id}
						/>
					</ListItem>
				)
			);

		var sourceLists = this.state.sourceLists.map(
			(l, i, ls) => (
				<ListItem key={l.id} list={l}>
					<RemovableItemFooter
						onRemove={this.onRemoveFromSource.bind(this, i)}
					/>
				</ListItem>
			)
		);
		if (sourceLists.length === 0) {
			sourceLists = <p><em>No source lists selected</em></p>;
		}

		var destList = <p><em>No destination list selected</em></p>;
		if (this.state.destList) {
			destList = (
				<ListItem list={this.state.destList}>
					<RemovableItemFooter
						onRemove={this.onRemoveDest.bind(this)}
					/>
				</ListItem>
			);
		}

		return (
			<div>
				<div className="left_column">
					<h2>Available Playlists</h2>
					{freeLists}
				</div>
				<div className="middle_column">
					<h2>Source Lists</h2>
					{sourceLists}
					<h2>Destination List</h2>
					{destList}
				</div>
				<div className="right_column">
					<h2>Some stuff!</h2>
				</div>
			</div>
		);
	}
}
Composer.propTypes = {
	token: React.PropTypes.string.isRequired,
	userID: React.PropTypes.string,
	playlists: React.PropTypes.array,
};

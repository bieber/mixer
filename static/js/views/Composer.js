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
			roundRobin: false,
			shuffle: false,
			dedup: false,
			pad: false,
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

	onFlipBool(field, event) {
		var update = {};
		update[field] = event.target.checked;
		this.setState(update);
	}

	onSubmit(event) {
		event.preventDefault();

		var {
			sourceLists,
			destList,
			roundRobin,
			shuffle,
			dedup,
			pad,
		} = this.state;

		function listExtractor(list) {
			return {
				id: list.id,
				owner_id: list.owner.id,
			};
		}

		this.props.onSubmit({
			source_lists: sourceLists.map((l, i, ls) => listExtractor(l)),
			dest_list: listExtractor(destList),
			options: {
				round_robin: roundRobin,
				shuffle: shuffle,
				dedup: dedup,
				pad: pad,
			},
		});
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

		var isDestEligible = l => {
			return this.props.userID === l.owner.id || l.collaborative;
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
							isDestEligible={isDestEligible(l)}
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

		var readyToSubmit = this.state.sourceLists.length > 0
			&& this.state.destList !== null;

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
					<h2>Mix</h2>
					<label>
						<input
							type="checkbox"
							checked={this.state.roundRobin}
							onChange={this.onFlipBool.bind(this, 'roundRobin')}
						/>
						<span>
							<strong>Round Robin</strong>
							<br />
							Insert one song at a time from each playlist.
						</span>
					</label>
					<br /><br />
					<label>
						<input
							type="checkbox"
							checked={this.state.shuffle}
							onChange={this.onFlipBool.bind(this, 'shuffle')}
						/>
						<span>
							<strong>Shuffle</strong>
							<br />
							Insert songs in random order.
						</span>
					</label>
					<br /><br />
					<label>
						<input
							type="checkbox"
							checked={this.state.dedup}
							onChange={this.onFlipBool.bind(this, 'dedup')}
						/>
						<span>
							<strong>Deduplicate</strong>
							<br />
							If a song appears in multiple source playlists, only
							insert it once into the output.
						</span>
					</label>
					<br /><br />
					<label>
						<input
							type="checkbox"
							checked={this.state.pad}
							onChange={this.onFlipBool.bind(this, 'pad')}
						/>
						<span>
							<strong>Pad</strong>
							<br />
							Repeat songs from shorter playlists to equal the
							length of longer playlists.
						</span>
					</label>
					<br /><br />
					<input
						type="submit"
						value="Mix"
						disabled={!readyToSubmit}
						onClick={this.onSubmit.bind(this)}
					/>
				</div>
			</div>
		);
	}
}
Composer.propTypes = {
	token: React.PropTypes.string.isRequired,
	userID: React.PropTypes.string,
	playlists: React.PropTypes.array,
	onSubmit: React.PropTypes.func.isRequired,
};

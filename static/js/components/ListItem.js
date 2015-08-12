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

import React from 'react';

export default class ListItem extends React.Component {
	render() {
		var properties = [];
		properties.push(this.props.list.public ? 'Public' : 'Private');
		if (this.props.list.collaborative) {
			properties.push('Collaborative');
		}
		
		return (
			<div className="list_item">
				<h3>
					{this.props.list.name}
				</h3>
				<p>by <strong>{this.props.list.owner.id}</strong></p>
				<p>{properties.join(', ')}</p>
				{this.props.children}
			</div>
		);
	}
}
ListItem.propTypes = {
	list: React.PropTypes.shape({
		id: React.PropTypes.string.isRequired,
		name: React.PropTypes.string.isRequired,
		collaborative: React.PropTypes.bool.isRequired,
		public: React.PropTypes.bool.isRequired,
		owner: React.PropTypes.shape({
			id: React.PropTypes.string,
		}),
	}),
};

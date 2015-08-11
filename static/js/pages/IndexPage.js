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

import Intro from '../views/Intro.js';

export default class IndexPage extends React.Component {
	constructor(props, context) {
		super(props, context);
		this.state = {
			someStuff: null,
		};
	}

	onLogin(data) {
		console.log(data);
	}

	render() {
		var view = (
			<Intro
				loginURI={this.props.loginURI}
				onLogin={this.onLogin.bind(this)}
			/>
		);

		return (
			<div className="container">
				{view}
			</div>
		);
	}
}
IndexPage.propTypes = {
	loginURI: React.PropTypes.string.isRequired,
};

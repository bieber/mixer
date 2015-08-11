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

export default class Intro extends React.Component {
	constructor(props, context) {
		super(props, context);
		this.popup = null;
	}

	onLoginClick(event) {
		event.preventDefault();
		if (this.popup !== null && !this.popup.closed) {
			alert('You already have a login popup open');
			return;
		}

		var width = 400;
		var height = 600;

		var windowFeatures = {
			menubar: 'no',
			location: 'no',
			left: (window.screen.width - width) / 2,
			top: (window.screen.height - height) / 2,
			width: width,
			height: height,
		};

		var featureStrings = [];
		for (var i in windowFeatures) {
			featureStrings.push(i+'='+windowFeatures[i]);
		}

		window.onLogin = this.props.onLogin;
		this.popup = window.open(
			this.props.loginURI,
			'login_window',
			featureStrings.join(',')
		);
	}

	render() {
		return (
			<div>
				<h1>Playlist Mixer</h1>
				<p>
					You can use this mixer to combine multiple spotify
					playlists into one using different merge strategies to
					order songs coming from each list and/or pad out the list
					so each playlist's songs have an equal chance of coming up
					on shuffle play.
				</p>
				<p>
					Before you get started, there's some things I'm supposed
					to tell you.
				</p>
				<ul>
					<li>
						This application is not developed or endorsed by
						Spotify.
					</li>
					<li>
						To operate, it will need access to both read and write
						from your playlistst.  In particular, it will write to
						the playlist you select as a destination and destroy
						anything that's already there.  Don't set any playlist
						as a destination if you want to preserve the data in
						it.
					</li>
					<li>
						This application will not retain any data about you
						once you leave the page.  This may change in the future
						to facilitate automated synchronization, but for now
						it operates on a strictly manual basis.
					</li>
				</ul>
				<p>
					If you're okay with all of that, then go ahead and...
				</p>
				<div className="login_button_container">
					<a
						href="#"
						onClick={this.onLoginClick.bind(this)}>
						<img src="/static/img/login_button.png" />
					</a>
				</div>
			</div>
		);
	}
}
Intro.propTypes = {
	loginURI: React.PropTypes.string.isRequired,
	onLogin: React.PropTypes.func.isRequired,
};

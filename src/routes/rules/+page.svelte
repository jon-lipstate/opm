<div>
	Odin strives to ensure that programming is simple and straight-forward, as such the following rules are in place to
	try to keep the quality as high as possible.
</div>

<p>Registry Rules:</p>
<ul>
	<li>Max 25 dependencies including subtrees</li>
	<li>Must be a github user for 30 days per library submitted. (request exemption on discord)</li>
	<li>
		Style Guidance:
		<ul>
			<li>Libraries should be idiomatic odin</li>
			<li>
				Allocating procudures either have code-docs stating such and/or use as an argument: <pre
					style="display:inline">allocator:=context.allocator</pre>
			</li>
		</ul>
	</li>
	<li>Poor quality code may be removed from the registry, repeat takedowns result in loss of publishing priveleges</li>
	<li>Packages with 0.0.x semvers are considered prototypes and largely exempt from scrutiny</li>
	<li>
		Packages which are nominally empty or have little contents may be pruned by moderators after 90 days of no commits
	</li>
	<li>malware or hostile code will be removed and the account locked</li>
	<li>
		All packages require licenses to be submitted. If the package file mismatches the actual license file, the more
		permissive license shall be granted to users. You agree to this term by uploading to the registry.
	</li>

	<li>
		Activity: TBD, we may send emails after 180 days (1 year?) of inactivity to affirm permanently stable, eg "done",
		but maintained if there are bugs, or abandoned. No response after 7 days automatically marks the repository
		abandoned. This can be undone in the manage panel.
		<ul>
			<li>
				Archival/Abandonment: If you decide to abandon a project, you must mark it as such in the Manage tab, and/or
				archive it on github. This allows users to know that they are not going to see any support up front.
			</li>
		</ul>
	</li>
</ul>
